package net

import "xtnet/net/packet"

type OnAccept func(ISession)
type OnSessionData func(ISession, []byte)
type OnSessionPacket func(ISession, *packet.ReadPacket)
type OnSessionClose func(ISession)

type OnRpcDirect func(ISession, *packet.ReadPacket)
type OnRpcRequestAsyn func(ISession, int32, *packet.ReadPacket)
type OnRpcRequestSync func(ISession, int32, *packet.ReadPacket)

type ISession interface {
	Send([]byte)
}

type IAgent interface {
	HandlerAccept(ISession)
	HandlerSessionClose(ISession)
	HandlerSessionData(ISession, []byte)
}
