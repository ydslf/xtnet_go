package rpc

import (
	"errors"
	"sync"
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

const RequestTimeout = 5000

type ContextSync struct {
	contextID int32
	t         RequestType
	cb        RequestCallback
}

type Sync struct {
	loop         *frame.Loop
	onRpcDirect  OnRpcDirect
	onRpcRequest OnRpcRequest
	contextID    int32
	contexts     sync.Map
	datas        chan *packet.ReadPacket
}

func NewSync(loop *frame.Loop) IRpc {
	return &Sync{
		loop:  loop,
		datas: make(chan *packet.ReadPacket, 1),
	}
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
	if c, ok := rpc.contexts.Load(contextID); ok {
		rpc.contexts.Delete(contextID)
		context := c.(*ContextSync)
		if context.t == rqtAsync {
			rpc.loop.Post(func() {
				context.cb(rpk)
			})
		} else {
			rpc.datas <- rpk
		}
	} else {
		xtnet.GetLogger().LogWarn("rpc.handlerResponse: no context, contextID=%d", contextID)
	}
}

func (rpc *Sync) GenContextID() int32 {
	if rpc.contextID == maxContextID {
		rpc.contextID = 0
	}
	rpc.contextID++
	return rpc.contextID
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

func (rpc *Sync) RequestSync(session net.ISession, wpk *packet.WritePacket, expireMS int) (*packet.ReadPacket, error) {
	contextID := rpc.GenContextID()
	context := &ContextSync{
		contextID: contextID,
		t:         rqtSync,
	}
	rpc.contexts.Store(contextID, context)

	wpk.WriteReserveInt32(contextID)
	wpk.WriteReserveInt8(rtRequest)
	session.Send(wpk.GetRealData())

	select {
	case rpk, ok := <-rpc.datas:
		if !ok {
			return nil, errors.New("RequestSync rpc closed")
		}
		return rpk, nil
	case <-time.After(RequestTimeout):
		rpc.contexts.Delete(contextID)
		return nil, errors.New("RequestSync timeout")
	}
}

func (rpc *Sync) Respond(session net.ISession, contextID int32, wpk *packet.WritePacket) {
	wpk.WriteReserveInt32(contextID)
	wpk.WriteReserveInt8(rtResponse)
	session.Send(wpk.GetRealData())
}
