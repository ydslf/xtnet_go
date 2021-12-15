package rpc

import (
	"errors"
	xtnet "xtnet"
	"xtnet/frame"
	"xtnet/net"
	"xtnet/net/packet"
)

type ContextNoSync struct {
	contextID int32
	cb        RequestCallback
}

type NoSync struct {
	loop         *frame.Loop
	onRpcDirect  OnRpcDirect
	onRpcRequest OnRpcRequest
	contextID    int32
	contexts     map[int32]*ContextNoSync
}

func NewNoSync(loop *frame.Loop) IRpc {
	return &NoSync{
		loop:     loop,
		contexts: make(map[int32]*ContextNoSync),
	}
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
			context.cb(rpk)
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

func (rpc *NoSync) SendDirect(session net.ISession, wpk *packet.WritePacket) {
	if session != nil {
		wpk.WriteReserveInt32(0)
		wpk.WriteReserveInt8(rtDirect)
		session.Send(wpk.GetRealData())
	}
}

func (rpc *NoSync) RequestAsync(session net.ISession, wpk *packet.WritePacket, cb RequestCallback) {
	if session != nil {
		contextID := rpc.GenContextID()
		context := &ContextNoSync{
			contextID: contextID,
			cb:        cb,
		}
		rpc.contexts[contextID] = context

		wpk.WriteReserveInt32(contextID)
		wpk.WriteReserveInt8(rtRequest)
		session.Send(wpk.GetRealData())
	}
}

func (rpc *NoSync) RequestSync(session net.ISession, wpk *packet.WritePacket, expireMS int) (rpk *packet.ReadPacket, err error) {
	return nil, errors.New("this rpc do not support RequestSync")
}

func (rpc *NoSync) Respond(session net.ISession, contextID int32, wpk *packet.WritePacket) {
	if session != nil {
		wpk.WriteReserveInt32(contextID)
		wpk.WriteReserveInt8(rtResponse)
		session.Send(wpk.GetRealData())
	}
}
