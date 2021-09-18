package timer

import (
	_ "runtime"
	"time"
	_ "unsafe" // for go:linkname
	"xtnet/frame"
)

//go:linkname runtimeNano runtime.nanotime
func runtimeNano() int64

//go:linkname deltimer runtime.deltimer
func deltimer(*runtimeTimer) bool

//go:linkname addtimer runtime.addtimer
func addtimer(t *runtimeTimer)

func when(d time.Duration) int64 {
	if d <= 0 {
		return runtimeNano()
	}
	t := runtimeNano() + int64(d)
	if t < 0 {
		t = 1<<63 - 1 // math.MaxInt64
	}
	return t
}

type runtimeTimer struct {
	pp       uintptr
	when     int64
	period   int64
	f        func(interface{}, uintptr) // NOTE: must not be closure
	arg      interface{}
	seq      uintptr
	nextwhen int64
	status   uint32
}

func systemTimerFunc(arg interface{}, seq uintptr) {
	timer := arg.(*SystemTimer)
	timer.loop.Post(func() {
		timer.cb()
	})
}

type SystemTimer struct {
	loop *frame.Loop
	cb   Cb
	r    runtimeTimer
}

func NewSystemTimer(loop *frame.Loop) *SystemTimer {
	return &SystemTimer{
		loop: loop,
	}
}

func (timer *SystemTimer) Start(d time.Duration, repeat bool, cb Cb) {
	//TODO 要先关闭
	timer.cb = cb
	timer.r.when = when(d)
	timer.r.f = systemTimerFunc
	timer.r.arg = timer
	if repeat {
		timer.r.period = int64(d)
	}

	addtimer(&timer.r)
}

func (timer *SystemTimer) Stop() {
	if timer.r.f == nil {
		panic("time: Stop called on uninitialized Timer")
	}
	deltimer(&timer.r)
}
