package timer

import (
	"sync"
	"time"
	"xtnet/frame"
)

type SystemTimer struct {
	loop     *frame.Loop
	ticker   *time.Ticker
	running  bool
	stopChan chan struct{}
	mutex    sync.Mutex
	wg       sync.WaitGroup
}

func NewSystemTimer(loop *frame.Loop) *SystemTimer {
	return &SystemTimer{
		loop:    loop,
		running: false,
	}
}

func (t *SystemTimer) Start(d time.Duration, repeat int, cb Cb) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if t.running {
		t.stop()
	}

	t.ticker = time.NewTicker(d)
	t.stopChan = make(chan struct{})
	t.running = true
	t.wg.Add(1)

	go t.run(repeat, cb)
}

func (t *SystemTimer) run(repeat int, cb Cb) {
	defer t.wg.Done()

	defer func() {
		t.mutex.Lock()
		defer t.mutex.Unlock()

		if t.ticker != nil {
			t.ticker.Stop()
			t.ticker = nil
		}
		t.stopChan = nil
		t.running = false
	}()

	count := 0
	for {
		select {
		case <-t.ticker.C:
			t.loop.Post(func() {
				cb()
			})

			if repeat == RepeatInfinity {
				continue
			}

			count++
			if count > repeat {
				return
			}

		case <-t.stopChan:
			return
		}
	}
}

func (t *SystemTimer) Stop() {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.stop()
}

func (t *SystemTimer) stop() {
	if !t.running {
		return
	}
	t.running = false

	if t.stopChan != nil {
		close(t.stopChan)
	}

	t.wg.Wait()
}
