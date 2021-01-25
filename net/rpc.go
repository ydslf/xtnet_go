package net

import (
	"xtnet/frame"
	"xtnet/net/packet"
)

const (
	RPCTYPE_DIRECT        int8 = 0 //直接，没有RPC
	RPCTYPE_REQUEST_ASYN  int8 = 1 //异步RPC请求
	PRCTYPE_REQUEST_SYNC  int8 = 2 //同步RPC请求
	RPCTYPE_RESPONSE_ASYN int8 = 3 //异步RPC回应
	RPCTYPE_RESPONSE_SYNC int8 = 4 //同步RPC回应
)

type RpcRequest struct {
	Session   ISession
	rpcType   int8
	contextID int32
}

type Rpc struct {
	loop         *frame.Loop
	onRpcDirect  OnRpcDirect
	onRpcRequest OnRpcRequest
}

func NewRpc(loop *frame.Loop) *Rpc {
	return &Rpc{
		loop: loop,
	}
}

func (rpc *Rpc) SetOnRpcDirect(onRpcDirect OnRpcDirect) {
	rpc.onRpcDirect = onRpcDirect
}

func (rpc *Rpc) SetOnRpcRequest(nRpcRequest OnRpcRequest) {
	rpc.onRpcRequest = nRpcRequest
}

func (rpc *Rpc) HandleSessionPacket(session ISession, rpk *packet.ReadPacket) {
	rpcType := rpk.ReadInt8()
	contextID := rpk.ReadInt32()

	switch rpcType {
	case RPCTYPE_DIRECT:
		rpc.loop.Post(func() {
			rpc.handleRpcDirect(session, rpk)
		})
	case RPCTYPE_REQUEST_ASYN:
		rpc.loop.Post(func() {
			rpc.handleRpcRequest(session, rpcType, contextID, rpk)
		})
	case PRCTYPE_REQUEST_SYNC:
		rpc.loop.Post(func() {
			rpc.handleRpcRequest(session, rpcType, contextID, rpk)
		})
	case RPCTYPE_RESPONSE_ASYN:
	case RPCTYPE_RESPONSE_SYNC:
	default:
	}
}

func (rpc *Rpc) handleRpcDirect(session ISession, rpk *packet.ReadPacket) {
	rpc.onRpcDirect(session, rpk)
}

func (rpc *Rpc) handleRpcRequest(session ISession, rpcType int8, contextID int32, rpk *packet.ReadPacket) {
	rpcRequest := &RpcRequest{
		Session:   session,
		rpcType:   rpcType,
		contextID: contextID,
	}
	rpc.onRpcRequest(rpcRequest, rpk)
}

func (rpc *Rpc) SendDirect(session ISession, wpk *packet.WritePacket) {

}

func (rpc *Rpc) RequestAsyn(session ISession, wpk *packet.WritePacket) {

}

func (rpc *Rpc) RequestSync(session ISession, wpk *packet.WritePacket) {

}

func (rpc *Rpc) Respond(rpcRequest *RpcRequest, wpk *packet.WritePacket) {

}
