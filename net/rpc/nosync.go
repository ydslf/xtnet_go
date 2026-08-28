package rpc

import (
	"errors"
	"sync"
	"time"
	xtnet "xtnet"
	"xtnet/frame"
	"xtnet/net"
	"xtnet/net/packet"
	xttimer "xtnet/timer"
)

type ContextNoSync struct {
	contextID int32
	cb        RequestCallback
	timer     *xttimer.WheelTimer
}

type NoSync struct {
	loop          *frame.Loop
	onRpcDirect   OnRpcDirect
	onRpcRequest  OnRpcRequest
	contextID     int32
	contexts      map[int32]*ContextNoSync
	timeWheelOnce sync.Once
	timeWheel     *xttimer.TimeWheel
}

func NewNoSync(loop *frame.Loop) IRpc {
	return &NoSync{
		loop:     loop,
		contexts: make(map[int32]*ContextNoSync),
	}
}

func (rpc *NoSync) newWheelTimer() *xttimer.WheelTimer {
	rpc.timeWheelOnce.Do(func() {
		rpc.timeWheel = xttimer.NewTimeWheel(rpc.loop, 0)
	})
	return rpc.timeWheel.NewTimer()
}

func (rpc *NoSync) SetOnRpcDirect(onRpcDirect OnRpcDirect) {
	rpc.onRpcDirect = onRpcDirect
}

func (rpc *NoSync) SetOnRpcRequest(nRpcRequest OnRpcRequest) {
	rpc.onRpcRequest = nRpcRequest
}

func (rpc *NoSync) HandleSessionPacket(session net.ISession, rpk *packet.ReadPacket) {
	rpcType := rpk.ReadInt8()
	contextID := rpk.ReadInt32()

	switch rpcType {
	case rtDirect:
		rpc.handleRpcDirect(session, rpk)
	case rtRequest:
		rpc.handleRpcRequest(session, contextID, rpk)
	case rtResponse:
		rpc.handlerResponse(contextID, rpk)
	default:
	}
}

func (rpc *NoSync) handleRpcDirect(session net.ISession, rpk *packet.ReadPacket) {
	rpc.loop.Post(func() {
		rpc.onRpcDirect(session, rpk)
	})
}

func (rpc *NoSync) handleRpcRequest(session net.ISession, contextID int32, rpk *packet.ReadPacket) {
	rpc.loop.Post(func() {
		rpc.onRpcRequest(session, contextID, rpk)
	})
}

func (rpc *NoSync) handlerResponse(contextID int32, rpk *packet.ReadPacket) {
	rpc.loop.Post(func() {
		if context, ok := rpc.contexts[contextID]; ok {
			delete(rpc.contexts, contextID)
			if context.timer != nil {
				context.timer.Stop()
			}
			context.cb(rpk, nil)
		} else {
			xtnet.GetLogger().LogWarn("rpc.handlerResponse: no context, contextID=%d", contextID)
		}
	})
}

func (rpc *NoSync) GenContextID() int32 {
	if rpc.contextID == maxContextID {
		rpc.contextID = 0
	}
	rpc.contextID++
	return rpc.contextID
}

func (rpc *NoSync) WriteDirectHead(wpk *packet.WritePacket) {
	wpk.WriteReserveInt32(0)
	wpk.WriteReserveInt8(rtDirect)
}

func (rpc *NoSync) SendDirect(session net.ISession, wpk *packet.WritePacket) {
	if session != nil {
		wpk.WriteReserveInt32(0)
		wpk.WriteReserveInt8(rtDirect)
		session.Send(wpk.GetRealData())
	}
}

func (rpc *NoSync) SendDirectRaw(session net.ISession, wpk *packet.WritePacket) {
	if session != nil {
		session.Send(wpk.GetRealData())
	}
}

// RequestAsync rpc.contexts 没有加锁，只支持在 rpc.loop 循环中调用
func (rpc *NoSync) RequestAsync(session net.ISession, wpk *packet.WritePacket, expireMS time.Duration, cb RequestCallback) {
	if session != nil {
		if expireMS < minRequestTimeout {
			expireMS = minRequestTimeout
		}

		contextID := rpc.GenContextID()
		context := &ContextNoSync{
			contextID: contextID,
			cb:        cb,
			timer:     rpc.newWheelTimer(),
		}
		rpc.contexts[contextID] = context
		context.timer.Start(expireMS, 0, func() {
			if current, ok := rpc.contexts[contextID]; ok && current == context {
				delete(rpc.contexts, contextID)
				xtnet.GetLogger().LogWarn("rpc.RequestAsync: timeout, contextID=%d", contextID)
				context.cb(nil, ErrRequestTimeout)
			}
		})

		wpk.WriteReserveInt32(contextID)
		wpk.WriteReserveInt8(rtRequest)
		session.Send(wpk.GetRealData())
	}
}

func (rpc *NoSync) RequestSync(session net.ISession, wpk *packet.WritePacket, expireMS time.Duration) (rpk *packet.ReadPacket, err error) {
	return nil, errors.New("this rpc do not support RequestSync")
}

func (rpc *NoSync) Respond(session net.ISession, contextID int32, wpk *packet.WritePacket) {
	if session != nil {
		wpk.WriteReserveInt32(contextID)
		wpk.WriteReserveInt8(rtResponse)
		session.Send(wpk.GetRealData())
	}
}
