package net

import (
	"xtnet/net/packet"
)

type OnAccept func(ISession)
type OnConnected func(ISession)
type OnSessionData func(ISession, []byte)
type OnSessionPacket func(ISession, *packet.ReadPacket)
type OnSessionClose func(ISession)

type RpcRequest struct {
	Session   ISession
	RpcType   int8
	ContextID int32
}

type OnRpcDirect func(ISession, *packet.ReadPacket)
type OnRpcRequest func(*RpcRequest, *packet.ReadPacket)

type INetBase interface {
	OnSessionStarted(ISession)
	Start() bool
	Close()
}

type ISession interface {
	SetAgent(IAgent)
	Send([]byte)
	Close(waitWrite bool)
	CloseBlock(waitWrite bool)
}

type IAgent interface {
	HandlerAccept(ISession)
	HandlerConnected(ISession)
	HandlerSessionClose(ISession)
	HandlerSessionData(ISession, []byte)
}
