package frame

import (
	"sync"
)

type SubService struct {
	loop      *Loop
	loopReply *Loop
	wgClose   sync.WaitGroup
}

func NewSubService(loopSize int, loopFullWarn bool, loopReply *Loop) *SubService {
	return &SubService{
		loop:      NewLoop(loopSize, loopFullWarn),
		loopReply: loopReply,
	}
}

func (ss *SubService) Call(call LoopFun) {
	ss.loop.Post(call)
}

func (ss *SubService) Reply(reply LoopFun) {
	ss.loopReply.Post(reply)
}

func (ss *SubService) Run() {
	ss.wgClose.Add(1)
	go func() {
		defer ss.wgClose.Done()
		ss.loop.Run()
	}()
}

func (ss *SubService) Close(waitHandle bool) {
	ss.loop.Close(waitHandle)
	ss.wgClose.Wait()
}
