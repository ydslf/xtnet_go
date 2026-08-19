package frame

import (
	"runtime/debug"
	"sync"
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
	functions    []LoopFun
	running      []LoopFun
	functionHead int
	mutex        sync.Mutex
	wakeChan     chan struct{}
	closeChan    chan struct{}
	status       atomic.Int32
	waitHandle   bool
	fullWarn     bool
}

func NewLoop(size int, fullWarn bool) *Loop {
	if size < LoopSizeMin {
		size = LoopSizeMin
	}
	return &Loop{
		functions: make([]LoopFun, 0, size),
		wakeChan:  make(chan struct{}, 1),
		closeChan: make(chan struct{}),
		fullWarn:  fullWarn,
	}
}

func (loop *Loop) Post(f LoopFun) {
	status := loop.status.Load()
	if status == loopStatusClose {
		xtnet.GetLogger().LogError("Loop.Post: loop status=%d", status)
		return
	}

	loop.mutex.Lock()
	status = loop.status.Load()
	if status == loopStatusClose {
		loop.mutex.Unlock()
		xtnet.GetLogger().LogError("Loop.Post: loop status=%d", status)
		return
	}

	loop.compact()
	oldCap := cap(loop.functions)
	loop.functions = append(loop.functions, f)
	queueSize := len(loop.functions) - loop.functionHead
	newCap := cap(loop.functions)
	loop.mutex.Unlock()

	if loop.fullWarn && newCap > oldCap {
		xtnet.GetLogger().LogWarn("Loop.Post: queue grow size=%d cap=%d->%d", queueSize, oldCap, newCap)
	}

	select {
	case loop.wakeChan <- struct{}{}:
	default:
	}
}

func (loop *Loop) compact() {
	if loop.functionHead == 0 || len(loop.functions) < cap(loop.functions) {
		return
	}

	n := copy(loop.functions, loop.functions[loop.functionHead:])
	clear(loop.functions[n:])
	loop.functions = loop.functions[:n]
	loop.functionHead = 0
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
		case <-loop.wakeChan:
			if !loop.runQueued() {
				return
			}
		case <-loop.closeChan:
			if loop.waitHandle {
				loop.drainQueued()
			}
			return
		}
	}
}

func (loop *Loop) runQueued() bool {
	functions := loop.takeAll()
	for _, f := range functions {
		select {
		case <-loop.closeChan:
			if !loop.waitHandle {
				clear(functions)
				return false
			}
		default:
		}

		loop.protectFun(f)
	}
	clear(functions)
	return true
}

func (loop *Loop) drainQueued() {
	for {
		functions := loop.takeAll()
		if len(functions) == 0 {
			return
		}

		for _, f := range functions {
			loop.protectFun(f)
		}
		clear(functions)
	}
}

func (loop *Loop) takeAll() []LoopFun {
	loop.mutex.Lock()
	queued := loop.functions[loop.functionHead:]
	loop.functions, loop.running = loop.running[:0], loop.functions
	loop.functionHead = 0
	loop.mutex.Unlock()
	return queued
}

func (loop *Loop) pop() (LoopFun, bool) {
	loop.mutex.Lock()
	if loop.functionHead == len(loop.functions) {
		loop.mutex.Unlock()
		return nil, false
	}

	f := loop.functions[loop.functionHead]
	loop.functions[loop.functionHead] = nil
	loop.functionHead++
	more := loop.functionHead < len(loop.functions)
	if !more {
		loop.functions = loop.functions[:0]
		loop.functionHead = 0
	}
	loop.mutex.Unlock()

	if more {
		select {
		case loop.wakeChan <- struct{}{}:
		default:
		}
	}
	return f, true
}

func (loop *Loop) RunOnce() {
	status := loop.status.Load()
	if status != loopStatusInit {
		xtnet.GetLogger().LogWarn("Loop.Run: loop status=%d", status)
		return
	}

	select {
	case <-loop.wakeChan:
		if f, ok := loop.pop(); ok {
			f()
		}
	case <-loop.closeChan:
	}
}

func (loop *Loop) Close(waitHandle bool) {
	loop.mutex.Lock()
	for {
		status := loop.status.Load()
		if status == loopStatusClose {
			loop.mutex.Unlock()
			xtnet.GetLogger().LogError("Loop.Close: loop status=%d", status)
			return
		}
		if loop.status.CompareAndSwap(status, loopStatusClose) {
			break
		}
	}

	loop.waitHandle = waitHandle
	close(loop.closeChan)
	loop.mutex.Unlock()
}
