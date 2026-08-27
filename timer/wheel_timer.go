package timer

import (
	"container/list"
	"sync"
	"sync/atomic"
	"time"
	"xtnet/frame"
)

const (
	defaultWheelTickTime = time.Millisecond * 10
)

var defaultWheelSlotSizes = []int{100, 60, 60, 24, 365}

// TimeWheelLevel 时间轮的一层，高层覆盖更长的 tick 跨度。
type TimeWheelLevel struct {
	size       int          // 本层槽数
	slotTicks  int64        // 本层每个槽覆盖的 tick 数
	totalTicks int64        // 本层总的tick数
	currentPos int          // 本层当前的位置
	slots      []*list.List // 槽数组，每个槽是一个任务链表
}

// TimeWheel 分层时间轮，单 tick 驱动，mutex 保护内部状态。
type TimeWheel struct {
	mu          sync.Mutex        // 保护 levels / currentTick 及所有任务链表
	loop        *frame.Loop       // 回调投递目标的事件循环
	tickTime    time.Duration     // 每次 tick 的时间间隔
	levels      []TimeWheelLevel  // 各层时间轮
	currentTick int64             // 当前累计 tick 数
	ticker      *time.Ticker      // 底层时钟源
	ready       []*wheelTimerTask // 已到期、等待投递到事件循环的任务
}

type cbWrap struct{ fn Cb }

// wheelTimerTask 单个定时任务，挂在某一层某个槽的链表中。
type wheelTimerTask struct {
	loop     *frame.Loop            // 回调需投递到的事件循环
	duration time.Duration          // 任务间隔
	repeat   int                    // 剩余重复次数，RepeatInfinity(-1) 表示无限
	cb       atomic.Pointer[cbWrap] // 到期回调
	dueTick  int64                  // 到期 tick
	slot     *list.List             // 当前所在槽（nil 表示已取出或未入槽）
	element  *list.Element          // 在槽链表中的元素句柄
	round    int                    // 还需在最高层完整循环的轮数，0 表示最后一轮
	cancel   atomic.Bool            // 是否已取消
}

func NewTimeWheel(loop *frame.Loop, tickTime time.Duration, slotSizes ...int) *TimeWheel {
	if tickTime <= 0 {
		tickTime = defaultWheelTickTime
	}
	if len(slotSizes) == 0 {
		slotSizes = defaultWheelSlotSizes
	}

	levelCount := len(slotSizes)
	levels := make([]TimeWheelLevel, levelCount)

	slotTicks := int64(1)
	for i := 0; i < levelCount; i++ {
		slots := make([]*list.List, slotSizes[i])
		for j := range slots {
			slots[j] = list.New()
		}
		levels[i] = TimeWheelLevel{
			size:       slotSizes[i],
			slotTicks:  slotTicks,
			totalTicks: slotTicks * int64(slotSizes[i]),
			slots:      slots,
		}
		slotTicks *= int64(slotSizes[i])
	}

	wheel := &TimeWheel{
		loop:     loop,
		tickTime: tickTime,
		levels:   levels,
		ticker:   time.NewTicker(tickTime),
	}

	go wheel.tickLoop()
	return wheel
}

func (wheel *TimeWheel) tickLoop() {
	for range wheel.ticker.C {
		wheel.mu.Lock()
		wheel.tickOnce()
		ready := wheel.takeReady()
		wheel.mu.Unlock()
		wheel.postReady(ready)
	}
}

// add locks internally; safe from any goroutine.
func (wheel *TimeWheel) add(task *wheelTimerTask) {
	wheel.mu.Lock()
	defer wheel.mu.Unlock()
	wheel.schedule(task)
}

// remove locks internally; safe from any goroutine.
func (wheel *TimeWheel) remove(task *wheelTimerTask) {
	if task == nil {
		return
	}
	wheel.mu.Lock()
	defer wheel.mu.Unlock()
	task.cancel.Store(true)
	wheel.removeTask(task)
}

// tickOnce, schedule, insertTask, processLevel, fire and removeTask all
// assume the caller holds wheel.mu.

func (wheel *TimeWheel) tickOnce() {
	wheel.currentTick++

	for level := 0; level < len(wheel.levels); level++ {
		if wheel.currentTick%wheel.levels[level].slotTicks != 0 {
			break
		}
		wheel.processLevel(level)
	}
}

func (wheel *TimeWheel) fire(task *wheelTimerTask) {
	if task.cancel.Load() {
		return
	}

	wheel.ready = append(wheel.ready, task)

	if task.cancel.Load() {
		return
	}

	if task.repeat == RepeatInfinity {
		wheel.schedule(task)
		return
	}

	if task.repeat > 0 {
		task.repeat--
		wheel.schedule(task)
	}
}

// takeReady assumes wheel.mu is held and transfers ownership of the batch.
func (wheel *TimeWheel) takeReady() []*wheelTimerTask {
	ready := wheel.ready
	wheel.ready = nil
	return ready
}

// postReady deliberately runs without wheel.mu held.
func (wheel *TimeWheel) postReady(ready []*wheelTimerTask) {
	for _, task := range ready {
		if task.cancel.Load() {
			continue
		}

		task.loop.Post(func() {
			if !task.cancel.Load() {
				if w := task.cb.Load(); w != nil {
					w.fn()
				}
			}
		})
	}
}

func (wheel *TimeWheel) schedule(task *wheelTimerTask) {
	if task == nil || task.cancel.Load() {
		return
	}

	ticks := int64((task.duration + wheel.tickTime - 1) / wheel.tickTime)
	if ticks <= 0 {
		ticks = 1
	}

	task.dueTick = wheel.currentTick + ticks
	wheel.insertTask(task)
}

func (wheel *TimeWheel) insertTask(task *wheelTimerTask) {
	if task == nil || task.cancel.Load() {
		return
	}

	remaining := task.dueTick - wheel.currentTick
	if remaining <= 0 {
		wheel.fire(task)
		return
	}

	level := 0
	find := false
	for i := range wheel.levels {
		if remaining < wheel.levels[i].totalTicks {
			level = i
			find = true
			break
		}
	}

	if find {
		lv := &wheel.levels[level]
		pos := int((task.dueTick / lv.slotTicks) % int64(lv.size))
		slot := lv.slots[pos]
		task.slot = slot
		task.element = slot.PushBack(task)
	} else {
		level = len(wheel.levels) - 1
		lv := &wheel.levels[level]
		round := int(remaining / lv.totalTicks)
		pos := int((task.dueTick / lv.slotTicks) % int64(lv.size))
		posCur := int((wheel.currentTick / lv.slotTicks) % int64(lv.size))
		if pos <= posCur {
			round--
		}
		slot := lv.slots[pos]
		task.round = round
		task.slot = slot
		task.element = slot.PushBack(task)
	}
}

func (wheel *TimeWheel) processLevel(level int) {
	lv := &wheel.levels[level]
	slotIndex := int((wheel.currentTick / lv.slotTicks) % int64(lv.size))
	slot := lv.slots[slotIndex]

	for element := slot.Front(); element != nil; {
		next := element.Next()
		task := element.Value.(*wheelTimerTask)
		slot.Remove(element)
		task.slot = nil
		task.element = nil

		if task.cancel.Load() {
			element = next
			continue
		}

		if task.round > 0 {
			task.round--
		} else {
			if task.dueTick > wheel.currentTick {
				wheel.insertTask(task)
			} else {
				wheel.fire(task)
			}
		}

		element = next
	}
}

func (wheel *TimeWheel) removeTask(task *wheelTimerTask) {
	if task == nil || task.slot == nil || task.element == nil {
		return
	}

	task.slot.Remove(task.element)
	task.slot = nil
	task.element = nil
}

func (wheel *TimeWheel) NewTimer() *WheelTimer {
	return &WheelTimer{
		wheel: wheel,
	}
}

// WheelTimer 面向用户的时间轮定时器，线程安全。
type WheelTimer struct {
	wheel *TimeWheel                     // 所属时间轮
	task  atomic.Pointer[wheelTimerTask] // 当前任务指针，nil 表示已停止
}

func (timer *WheelTimer) Start(d time.Duration, repeat int, cb Cb) {
	timer.tryStop()

	task := &wheelTimerTask{
		loop:     timer.wheel.loop,
		duration: d,
		repeat:   repeat,
	}
	task.cb.Store(&cbWrap{fn: cb})
	timer.task.Store(task)
	timer.wheel.add(task)
}

func (timer *WheelTimer) Stop() {
	timer.tryStop()
}

func (timer *WheelTimer) tryStop() {
	old := timer.task.Swap(nil)
	if old == nil {
		return
	}
	old.cancel.Store(true)
	old.cb.Store(nil)
}
