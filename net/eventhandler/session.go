package eventhandler

import (
	"xtnet/net"
	"xtnet/net/packet"
)

type OnAccept func(net.ISession)
type OnSessionPacket func(net.ISession, *packet.ReadPacket)
type OnSessionClose func(net.ISession)

type Session struct {
	OnAccept        OnAccept
	OnSessionPacket OnSessionPacket
	OnSessionClose  OnSessionClose
}

func NewSessionEventHandler() *Session {
	return &Session{}
}
