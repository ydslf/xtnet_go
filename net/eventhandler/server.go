package eventhandler

import (
	"xtnet/net"
	"xtnet/net/packet"
)

type OnAccept func(net.IServer, net.ISession)
type OnSessionPacket func(net.IServer, net.ISession, *packet.ReadPacket)
type OnSessionClose func(net.IServer, net.ISession)

type Server struct {
	OnAccept        OnAccept
	OnSessionPacket OnSessionPacket
	OnSessionClose  OnSessionClose
}

func NewServerEventHandler() *Server {
	return &Server{}
}
