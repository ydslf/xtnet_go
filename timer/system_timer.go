package timer

import (
	_ "runtime"
	"time"
	_ "unsafe" // for go:linkname
)

//go:linkname runtimeNano runtime.nanotime
func runtimeNano() int64

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

//go:linkname startTimer runtime.startTimer
func startTimer(*runtimeTimer)

//go:linkname addtimer runtime.addtimer
func addtimer(t *runtimeTimer)

func goFunc(arg interface{}, seq uintptr) {
	go arg.(func())()
}

func AfterFunc(d time.Duration, f func()) {
	r := runtimeTimer{
		when: when(d),
		f:    goFunc,
		arg:  f,
	}
	addtimer(&r)
}
