package rpc

import (
	"errors"
	"xtnet/frame"
	"xtnet/net"
	"xtnet/net/packet"
)

type Sync struct {
	loop         *frame.Loop
	onRpcDirect  OnRpcDirect
	onRpcRequest OnRpcRequest
	contextID    int32
}

func NewSync(service *frame.Service) IRpc {
	return &Sync{
		loop: service.GetLoop(),
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

}

func (rpc *Sync) GenContextID() int32 {
	if rpc.contextID == maxContextID {
		rpc.contextID = 0
	}
	rpc.contextID++
	return rpc.contextID
}

func (rpc *Sync) SendDirect(session net.ISession, wpk *packet.WritePacket) {
	wpk.WriteReserveInt32(0)
	wpk.WriteReserveInt8(rtDirect)
	session.Send(wpk.GetRealData())
}

func (rpc *Sync) RequestAsync(session net.ISession, wpk *packet.WritePacket, cb RequestCallback) {

}

func (rpc *Sync) RequestSync(session net.ISession, wpk *packet.WritePacket) (rpk *packet.ReadPacket, err error) {
	return nil, errors.New("this rpc do not support RequestSync")
}

func (rpc *Sync) Respond(session net.ISession, contextID int32, wpk *packet.WritePacket) {
	wpk.WriteReserveInt32(contextID)
	wpk.WriteReserveInt8(rtResponse)
	session.Send(wpk.GetRealData())
}
