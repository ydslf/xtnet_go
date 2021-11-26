package net

import (
	"xtnet/net/packet"
)

type RpcRequest struct {
	Session   ISession
	RpcType   int8
	ContextID int32
}

type OnRpcDirect func(ISession, *packet.ReadPacket)
type OnRpcRequest func(*RpcRequest, *packet.ReadPacket)

type INetBase interface {
	OnSessionStarted(ISession)
}

type IServer interface {
	Start() bool
	Close()
}

type IClient interface {
	Connect() bool
	ConnectSync() bool
	Close()
}

type ISession interface {
	SetAgent(ISessionAgent)
	Send([]byte)
	Close(waitWrite bool)
}

type ISessionAgent interface {
	HandlerAccept(ISession)
	HandlerSessionClose(ISession)
	HandlerSessionData(ISession, []byte)
}

type IClientAgent interface {
	HandlerConnected(IClient)
	HandlerDisconnected(IClient)
	HandlerClientData(ISession, []byte)
	HandlerConnectBreak(ISession)
}
