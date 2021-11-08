package rpc

import (
	"xtnet/frame"
	"xtnet/net"
	"xtnet/net/packet"
)

const (
	RpcTypeDirect        int8 = 0 //直接，没有RPC
	RpcTypeRequestAsync  int8 = 1 //异步RPC请求
	PrcTypeRequestSync   int8 = 2 //同步RPC请求
	RpcTypeResponseAsync int8 = 3 //异步RPC回应
	RpcTypeResponseSync  int8 = 4 //同步RPC回应
)

/*
	|                      pktBody                     |
	|  msgID |             msgBody                     |
	| msgDirection | msgID |          msgBody          |
	| rpcType | contextID  | msgType | msgID | msgBody |
*/

//TODO 应该改为IAgent的子类

type Rpc struct {
	loop         *frame.Loop
	onRpcDirect  net.OnRpcDirect
	onRpcRequest net.OnRpcRequest
}

func NewRpc(loop *frame.Loop) *Rpc {
	return &Rpc{
		loop: loop,
	}
}

func (rpc *Rpc) SetOnRpcDirect(onRpcDirect net.OnRpcDirect) {
	rpc.onRpcDirect = onRpcDirect
}

func (rpc *Rpc) SetOnRpcRequest(nRpcRequest net.OnRpcRequest) {
	rpc.onRpcRequest = nRpcRequest
}

func (rpc *Rpc) HandleSessionPacket(session net.ISession, rpk *packet.ReadPacket) {
	rpcType := rpk.ReadInt8()
	contextID := rpk.ReadInt32()

	switch rpcType {
	case RpcTypeDirect:
		rpc.loop.Post(func() {
			rpc.handleRpcDirect(session, rpk)
		})
	case RpcTypeRequestAsync:
		rpc.loop.Post(func() {
			rpc.handleRpcRequest(session, rpcType, contextID, rpk)
		})
	case PrcTypeRequestSync:
		rpc.loop.Post(func() {
			rpc.handleRpcRequest(session, rpcType, contextID, rpk)
		})
	case RpcTypeResponseAsync:
	case RpcTypeResponseSync:
	default:
	}
}

func (rpc *Rpc) handleRpcDirect(session net.ISession, rpk *packet.ReadPacket) {
	rpc.onRpcDirect(session, rpk)
}

func (rpc *Rpc) handleRpcRequest(session net.ISession, rpcType int8, contextID int32, rpk *packet.ReadPacket) {
	rpcRequest := &net.RpcRequest{
		Session:   session,
		RpcType:   rpcType,
		ContextID: contextID,
	}
	rpc.onRpcRequest(rpcRequest, rpk)
}

func (rpc *Rpc) SendDirect(session net.ISession, wpk *packet.WritePacket) {

}

func (rpc *Rpc) RequestAsync(session net.ISession, wpk *packet.WritePacket) {

}

func (rpc *Rpc) RequestSync(session net.ISession, wpk *packet.WritePacket) {

}

func (rpc *Rpc) Respond(rpcRequest *net.RpcRequest, wpk *packet.WritePacket) {

}
