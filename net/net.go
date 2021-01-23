package net

import "xtnet/net/packet"

type OnAccept func(Session)
type OnSessionData func(Session, []byte)
type OnSessionPacket func(Session, *packet.ReadPacket)
type OnSessionClose func(Session)

type OnRpcDirect func(Session, *packet.ReadPacket)
type OnRpcRequestAsyn func(Session, int32, *packet.ReadPacket)
type OnRpcRequestSync func(Session, int32, *packet.ReadPacket)

type Session interface {
	Send([]byte)
}

type IAgent interface {
	HandlerAccept(Session)
	HandlerSessionClose(Session)
	HandlerSessionData(Session, []byte)
}
