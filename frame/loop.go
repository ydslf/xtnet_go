package frame

import (
	"runtime/debug"
	"sync/atomic"
	xtnet "xtnet"
)

const (
	loopStatusInit int32 = iota
	loopStatusRunning
	loopStatusClose
)

const LoopSizeMin int = 4096

type LoopFun func()

type Loop struct {
	functions chan LoopFun
	closeChan chan int
	status    atomic.Int32
	fullWarn  bool
}

func NewLoop(size int, fullWarn bool) *Loop {
	if size < LoopSizeMin {
		size = LoopSizeMin
	}
	return &Loop{
		functions: make(chan LoopFun, size),
		closeChan: make(chan int, 1),
		fullWarn:  fullWarn,
	}
}

func (loop *Loop) Post(f LoopFun) {
	status := loop.status.Load()
	if status == loopStatusClose {
		xtnet.GetLogger().LogError("Loop.Post: loop status=%d", status)
		return
	}
	if loop.fullWarn {
		if len(loop.functions) > cap(loop.functions) {
			xtnet.GetLogger().LogWarn("Loop.Post: chan cap=%d", cap(loop.functions))
		}
	}

	loop.functions <- f
}

func (loop *Loop) protectFun(f LoopFun) {
	defer func() {
		if err := recover(); err != nil {
			xtnet.GetLogger().LogError("loop.protectFun: %v", err)
			xtnet.GetLogger().LogError(string(debug.Stack()))
		}
	}()
	f()
}

func (loop *Loop) Run() {
	if !loop.status.CompareAndSwap(loopStatusInit, loopStatusRunning) {
		xtnet.GetLogger().LogWarn("Loop.Run: loop status=%d", loop.status.Load())
		return
	}

	for {
		select {
		case f, ok := <-loop.functions:
			if ok == false {
				return
			}
			loop.protectFun(f)
		case <-loop.closeChan:
			return
		}
	}
}

func (loop *Loop) RunOnce() {
	status := loop.status.Load()
	if status != loopStatusInit {
		xtnet.GetLogger().LogWarn("Loop.Run: loop status=%d", status)
		return
	}

	f, ok := <-loop.functions
	if ok {
		f()
	}
}

func (loop *Loop) Close(waitHandle bool) {
	status := loop.status.Load()
	if status == loopStatusClose {
		xtnet.GetLogger().LogError("Loop.Close: loop status=%d", status)
		return
	}

	loop.status.Store(loopStatusClose)

	if waitHandle {
		close(loop.functions)
	} else {
		loop.closeChan <- 1
	}
}
