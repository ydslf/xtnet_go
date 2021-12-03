package rpc

import (
	"errors"
	xtnet "xtnet"
	"xtnet/frame"
	"xtnet/net"
	"xtnet/net/packet"
)

/*
	|                      pktBody                     |
	|  msgID |             msgBody                     |
	| rpcType | contextID  |       msgBody             |
*/

const (
	rtDirect   int8 = iota //直接，没有RPC
	rtRequest              //RPC请求
	rtResponse             //异步RPC回应
)

const maxContextID int32 = 0x8FFFFFFF

type OnRpcDirect func(session net.ISession, rpk *packet.ReadPacket)
type OnRpcRequest func(session net.ISession, contextID int32, rpk *packet.ReadPacket)
type RequestCallback func(rpk *packet.ReadPacket)

type Context struct {
	contextID int32
	cb        RequestCallback
}

type Rpc struct {
	loop         *frame.Loop
	onRpcDirect  OnRpcDirect
	onRpcRequest OnRpcRequest
	contextID    int32
	contexts     map[int32]*Context
}

func NewRpc(service *frame.Service) *Rpc {
	return &Rpc{
		loop:     service.GetLoop(),
		contexts: make(map[int32]*Context),
	}
}

func (rpc *Rpc) SetOnRpcDirect(onRpcDirect OnRpcDirect) {
	rpc.onRpcDirect = onRpcDirect
}

func (rpc *Rpc) SetOnRpcRequest(nRpcRequest OnRpcRequest) {
	rpc.onRpcRequest = nRpcRequest
}

func (rpc *Rpc) HandleSessionPacket(session net.ISession, rpk *packet.ReadPacket) {
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

func (rpc *Rpc) handleRpcDirect(session net.ISession, rpk *packet.ReadPacket) {
	rpc.loop.Post(func() {
		rpc.onRpcDirect(session, rpk)
	})
}

func (rpc *Rpc) handleRpcRequest(session net.ISession, contextID int32, rpk *packet.ReadPacket) {
	rpc.loop.Post(func() {
		rpc.onRpcRequest(session, contextID, rpk)
	})
}

func (rpc *Rpc) handlerResponse(contextID int32, rpk *packet.ReadPacket) {
	rpc.loop.Post(func() {
		context, ok := rpc.contexts[contextID]
		if ok {
			delete(rpc.contexts, contextID)
			context.cb(rpk)
		} else {
			xtnet.GetLogger().LogWarn("rpc.handlerResponse: no context, contextID=%d", contextID)
		}
	})
}

func (rpc *Rpc) GenContextID() int32 {
	if rpc.contextID == maxContextID {
		rpc.contextID = 0
	}
	rpc.contextID++
	return rpc.contextID
}

func (rpc *Rpc) SendDirect(session net.ISession, wpk *packet.WritePacket) {
	wpk.WriteReserveInt32(0)
	wpk.WriteReserveInt8(rtDirect)
	session.Send(wpk.GetRealData())
}

func (rpc *Rpc) RequestAsync(session net.ISession, wpk *packet.WritePacket, cb RequestCallback) {
	contextID := rpc.GenContextID()
	context := &Context{
		contextID: contextID,
		cb:        cb,
	}
	rpc.contexts[contextID] = context

	wpk.WriteReserveInt32(contextID)
	wpk.WriteReserveInt8(rtRequest)
	session.Send(wpk.GetRealData())
}

func (rpc *Rpc) RequestSync(session net.ISession, wpk *packet.WritePacket) (rpk *packet.ReadPacket, err error) {
	return nil, errors.New("this rpc do not support RequestSync")
}

func (rpc *Rpc) Respond(session net.ISession, contextID int32, wpk *packet.WritePacket) {
	wpk.WriteReserveInt32(contextID)
	wpk.WriteReserveInt8(rtResponse)
	session.Send(wpk.GetRealData())
}
