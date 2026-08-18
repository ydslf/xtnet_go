package frame

import (
	"runtime"
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
	functions  chan LoopFun
	closeChan  chan struct{}
	status     atomic.Int32
	posting    atomic.Int32
	waitHandle bool
	fullWarn   bool
}

func NewLoop(size int, fullWarn bool) *Loop {
	if size < LoopSizeMin {
		size = LoopSizeMin
	}
	return &Loop{
		functions: make(chan LoopFun, size),
		closeChan: make(chan struct{}),
		fullWarn:  fullWarn,
	}
}

func (loop *Loop) Post(f LoopFun) {
	loop.posting.Add(1)
	defer loop.posting.Add(-1)

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

	select {
	case loop.functions <- f:
	case <-loop.closeChan:
	}
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
		case f := <-loop.functions:
			loop.protectFun(f)
		case <-loop.closeChan:
			if loop.waitHandle {
				loop.drain()
			}
			return
		}
	}
}

func (loop *Loop) drain() {
	for {
		select {
		case f := <-loop.functions:
			loop.protectFun(f)
		default:
			if loop.posting.Load() == 0 {
				return
			}
			runtime.Gosched()
		}
	}
}

func (loop *Loop) RunOnce() {
	status := loop.status.Load()
	if status != loopStatusInit {
		xtnet.GetLogger().LogWarn("Loop.Run: loop status=%d", status)
		return
	}

	select {
	case f := <-loop.functions:
		f()
	case <-loop.closeChan:
	}
}

func (loop *Loop) Close(waitHandle bool) {
	for {
		status := loop.status.Load()
		if status == loopStatusClose {
			xtnet.GetLogger().LogError("Loop.Close: loop status=%d", status)
			return
		}
		if loop.status.CompareAndSwap(status, loopStatusClose) {
			break
		}
	}

	loop.waitHandle = waitHandle
	close(loop.closeChan)
}
