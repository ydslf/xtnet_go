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

type Rpc struct {
	loop             *frame.Loop
	onRpcDirect      OnRpcDirect
	onRpcRequestAsyn OnRpcRequestAsyn
	onRpcRequestSync OnRpcRequestSync
}

func NewRpc(loop *frame.Loop) *Rpc {
	return &Rpc{
		loop: loop,
	}
}

func (rpc *Rpc) SetOnRpcDirect(onRpcDirect OnRpcDirect) {
	rpc.onRpcDirect = onRpcDirect
}

func (rpc *Rpc) SetOnRpcRequestAsyn(onRpcRequestAsyn OnRpcRequestAsyn) {
	rpc.onRpcRequestAsyn = onRpcRequestAsyn
}

func (rpc *Rpc) SetOnRpcRequestSync(onRpcRequestSync OnRpcRequestSync) {
	rpc.onRpcRequestSync = onRpcRequestSync
}

func (rpc *Rpc) HandleSessionPacket(session Session, rpk *packet.ReadPacket) {
	rpcType := rpk.ReadInt8()
	contextID := rpk.ReadInt32()

	switch rpcType {
	case RPCTYPE_DIRECT:
		rpc.loop.Post(func() {
			rpc.handleRpcDirect(session, rpk)
		})
	case RPCTYPE_REQUEST_ASYN:
		rpc.loop.Post(func() {
			rpc.handleRpcRequestAsyn(session, contextID, rpk)
		})
	case PRCTYPE_REQUEST_SYNC:
		rpc.loop.Post(func() {
			rpc.handleRpcRequestSync(session, contextID, rpk)
		})
	case RPCTYPE_RESPONSE_ASYN:
	case RPCTYPE_RESPONSE_SYNC:
	default:
	}
}

func (rpc *Rpc) handleRpcDirect(session Session, rpk *packet.ReadPacket) {
	rpc.onRpcDirect(session, rpk)
}

func (rpc *Rpc) handleRpcRequestAsyn(session Session, contextID int32, rpk *packet.ReadPacket) {
	rpc.onRpcRequestAsyn(session, contextID, rpk)
}

func (rpc *Rpc) handleRpcRequestSync(session Session, contextID int32, rpk *packet.ReadPacket) {
	rpc.onRpcRequestSync(session, contextID, rpk)
}

func (rpc *Rpc) SendDirect() {

}

func (rpc *Rpc) RequestAsyn() {

}

func (rpc *Rpc) RequestSync() {

}

func (rpc *Rpc) RespondAsyn() {

}

func (rpc *Rpc) RespondSync() {

}
