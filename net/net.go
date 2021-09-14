package net

import "xtnet/net/packet"

type OnAccept func(ISession)
type OnSessionData func(ISession, []byte)
type OnSessionPacket func(ISession, *packet.ReadPacket)
type OnSessionClose func(ISession)

type OnRpcDirect func(ISession, *packet.ReadPacket)
type OnRpcRequest func(*RpcRequest, *packet.ReadPacket)

type ISession interface {
	Send([]byte)
	Close(waitWrite bool)
	CloseBlock(waitWrite bool)
}

type IAgent interface {
	HandlerAccept(ISession)
	HandlerSessionClose(ISession)
	HandlerSessionData(ISession, []byte)
}
