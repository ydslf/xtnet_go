package net

import (
	"xtnet/net/packet"
)

type RpcRequest struct {
	Session   ISession
	RpcType   int8
	ContextID int32
}

type OnSessionStart func(ISession)
type OnSessionData func(ISession, []byte)
type OnSessionClose func(ISession)

type OnRpcDirect func(ISession, *packet.ReadPacket)
type OnRpcRequest func(*RpcRequest, *packet.ReadPacket)

type IServer interface {
	Start() bool
	Close()
}

type IClient interface {
	Connect() bool
	ConnectSync(ms int) error
	Close(waitWrite bool)
}

type ISession interface {
	Send([]byte)
	Close(waitWrite bool)
	SetSessionStartCb(OnSessionStart)
	SetSessionDataCb(OnSessionData)
	SetSessionCloseCb(OnSessionClose)
}

type IServerAgent interface {
	HandlerAccept(IServer, ISession)
	HandlerSessionData(IServer, ISession, []byte)
	HandlerSessionClose(IServer, ISession)
}

type IClientAgent interface {
	HandlerConnect(IClient)
	HandlerDisconnect(IClient)
	HandlerClientData(IClient, []byte)
	HandlerConnectBreak(IClient)
}
