package rpc

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"
	xtnet "xtnet"
	"xtnet/frame"
	"xtnet/net"
	"xtnet/net/packet"
)

type RequestType int8

const (
	rqtAsync RequestType = iota
	rqtSync
)

const minRequestTimeout = 500 * time.Millisecond

type ContextSync struct {
	contextID int32
	t         RequestType
	cb        RequestCallback
	response  chan *packet.ReadPacket
}

type Sync struct {
	loop         *frame.Loop
	onRpcDirect  OnRpcDirect
	onRpcRequest OnRpcRequest
	contextID    atomic.Int32
	contexts     sync.Map
}

func NewSync(loop *frame.Loop) IRpc {
	return &Sync{loop: loop}
}

func (rpc *Sync) SetOnRpcDirect(onRpcDirect OnRpcDirect) {
	rpc.onRpcDirect = onRpcDirect
}

func (rpc *Sync) SetOnRpcRequest(nRpcRequest OnRpcRequest) {
	rpc.onRpcRequest = nRpcRequest
}

func (rpc *Sync) HandleSessionPacket(session net.ISession, rpk *packet.ReadPacket) {
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

func (rpc *Sync) handleRpcDirect(session net.ISession, rpk *packet.ReadPacket) {
	rpc.loop.Post(func() {
		rpc.onRpcDirect(session, rpk)
	})
}

func (rpc *Sync) handleRpcRequest(session net.ISession, contextID int32, rpk *packet.ReadPacket) {
	rpc.loop.Post(func() {
		rpc.onRpcRequest(session, contextID, rpk)
	})
}

func (rpc *Sync) handlerResponse(contextID int32, rpk *packet.ReadPacket) {
	if c, ok := rpc.contexts.LoadAndDelete(contextID); ok {
		context := c.(*ContextSync)
		if context.t == rqtAsync {
			rpc.loop.Post(func() {
				context.cb(rpk)
			})
		} else {
			context.response <- rpk
		}
	} else {
		xtnet.GetLogger().LogWarn("rpc.handlerResponse: no context, contextID=%d", contextID)
	}
}

func (rpc *Sync) GenContextID() int32 {
	for {
		current := rpc.contextID.Load()
		next := current + 1
		if current >= maxContextID {
			next = 1
		}
		if rpc.contextID.CompareAndSwap(current, next) {
			return next
		}
	}
}

func (rpc *Sync) WriteDirectHead(wpk *packet.WritePacket) {
	wpk.WriteReserveInt32(0)
	wpk.WriteReserveInt8(rtDirect)
}

func (rpc *Sync) SendDirect(session net.ISession, wpk *packet.WritePacket) {
	wpk.WriteReserveInt32(0)
	wpk.WriteReserveInt8(rtDirect)
	session.Send(wpk.GetRealData())
}

func (rpc *Sync) SendDirectRaw(session net.ISession, wpk *packet.WritePacket) {
	if session != nil {
		session.Send(wpk.GetRealData())
	}
}

func (rpc *Sync) RequestAsync(session net.ISession, wpk *packet.WritePacket, cb RequestCallback) {
	contextID := rpc.GenContextID()
	context := &ContextSync{
		contextID: contextID,
		t:         rqtAsync,
		cb:        cb,
	}
	rpc.contexts.Store(contextID, context)

	wpk.WriteReserveInt32(contextID)
	wpk.WriteReserveInt8(rtRequest)
	session.Send(wpk.GetRealData())
}

func (rpc *Sync) RequestSync(session net.ISession, wpk *packet.WritePacket, expireMS time.Duration) (*packet.ReadPacket, error) {
	if expireMS < minRequestTimeout {
		expireMS = minRequestTimeout
	}

	contextID := rpc.GenContextID()
	context := &ContextSync{
		contextID: contextID,
		t:         rqtSync,
		response:  make(chan *packet.ReadPacket, 1),
	}
	rpc.contexts.Store(contextID, context)

	wpk.WriteReserveInt32(contextID)
	wpk.WriteReserveInt8(rtRequest)
	session.Send(wpk.GetRealData())

	timer := time.NewTimer(expireMS)
	defer timer.Stop()

	select {
	case rpk := <-context.response:
		return rpk, nil
	case <-timer.C:
		rpc.contexts.Delete(contextID)
		return nil, errors.New("RequestSync timeout")
	}
}

func (rpc *Sync) Respond(session net.ISession, contextID int32, wpk *packet.WritePacket) {
	wpk.WriteReserveInt32(contextID)
	wpk.WriteReserveInt8(rtResponse)
	session.Send(wpk.GetRealData())
}
