package frame

import (
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
	status    int32
	warnLen   bool
}

func NewLoop(size int, warnLen bool) *Loop {
	if size < LoopSizeMin {
		size = LoopSizeMin
	}
	return &Loop{
		functions: make(chan LoopFun, size),
		closeChan: make(chan int, 1),
		status:    loopStatusInit,
		warnLen:   warnLen,
	}
}

func (loop *Loop) Post(f LoopFun) {
	if loop.status == loopStatusClose {
		xtnet.GetLogger().LogError("Loop.Post: loop status=%d", loop.status)
		return
	}
	if loop.warnLen {
		if len(loop.functions) > cap(loop.functions) {
			xtnet.GetLogger().LogWarn("Loop.Post: chan cap=%d", cap(loop.functions))
		}
	}

	loop.functions <- f
}

func (loop *Loop) Run() {
	if loop.status != loopStatusInit {
		xtnet.GetLogger().LogWarn("Loop.Run: loop status=%d", loop.status)
		return
	}

	loop.status = loopStatusRunning

	for {
		select {
		case f, ok := <-loop.functions:
			if ok == false {
				return
			}
			f()
		case <-loop.closeChan:
			return
		}
	}
}

func (loop *Loop) RunOnce() {
	if loop.status != loopStatusInit {
		xtnet.GetLogger().LogWarn("Loop.Run: loop status=%d", loop.status)
		return
	}

	f, ok := <-loop.functions
	if ok {
		f()
	}
}

func (loop *Loop) Close(waitHandle bool) {
	if loop.status == loopStatusClose {
		xtnet.GetLogger().LogError("Loop.Close: loop status=%d", loop.status)
		return
	}

	loop.status = loopStatusClose

	if waitHandle {
		close(loop.functions)
	} else {
		loop.closeChan <- 1
	}
}
