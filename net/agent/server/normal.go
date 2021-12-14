package server

import (
	"encoding/binary"
	"xtnet/frame"
	"xtnet/net"
	"xtnet/net/eventhandler"
	"xtnet/net/packet"
)

type Normal struct {
	loop         *frame.Loop
	byteOrder    binary.ByteOrder
	eventHandler *eventhandler.Server
}

func NewNormal(loop *frame.Loop, byteOrder binary.ByteOrder) *Normal {
	return &Normal{
		loop:      loop,
		byteOrder: byteOrder,
	}
}

func (agent *Normal) SetEventHandler(eventHandler *eventhandler.Server) {
	agent.eventHandler = eventHandler
}

func (agent *Normal) HandleAccept(server net.IServer, session net.ISession) {
	agent.loop.Post(func() {
		agent.eventHandler.OnAccept(server, session)
	})
}

func (agent *Normal) HandleSessionData(server net.IServer, session net.ISession, data []byte) {
	rpk := packet.NewReadPacket(data, agent.byteOrder, 0, len(data))
	agent.loop.Post(func() {
		agent.eventHandler.OnSessionPacket(server, session, rpk)
	})
}

func (agent *Normal) HandleSessionClose(server net.IServer, session net.ISession) {
	agent.loop.Post(func() {
		agent.eventHandler.OnSessionClose(server, session)
	})
}
