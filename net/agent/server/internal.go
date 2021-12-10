package server

import (
	"encoding/binary"
	"xtnet/frame"
	"xtnet/net"
	"xtnet/net/eventhandler"
	"xtnet/net/packet"
	"xtnet/net/rpc"
)

type Internal struct {
	loop         *frame.Loop
	eventHandler *eventhandler.Server
	netRpc       rpc.IRpc
}

func NewInternal(loop *frame.Loop) *Internal {
	return &Internal{
		loop: loop,
	}
}

func (agent *Internal) SetEventHandler(eventHandler *eventhandler.Server) {
	agent.eventHandler = eventHandler
}

func (agent *Internal) SetNetRpc(netRpc rpc.IRpc) {
	agent.netRpc = netRpc
}

func (agent *Internal) HandleAccept(server net.IServer, session net.ISession) {
	agent.loop.Post(func() {
		agent.eventHandler.OnAccept(server, session)
	})
}

func (agent *Internal) HandleSessionData(server net.IServer, session net.ISession, data []byte) {
	rpk := packet.NewReadPacket(data, binary.BigEndian, 0, len(data))
	agent.netRpc.HandleSessionPacket(session, rpk)
}

func (agent *Internal) HandleSessionClose(server net.IServer, session net.ISession) {
	agent.loop.Post(func() {
		agent.eventHandler.OnSessionClose(server, session)
	})
}
