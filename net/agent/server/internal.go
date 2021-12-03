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

func NewInternal(service *frame.Service) net.IServerAgent {
	return &Internal{
		loop: service.GetLoop(),
	}
}

func (agent *Internal) SetEventHandler(eventHandler *eventhandler.Server) {
	agent.eventHandler = eventHandler
}

func (agent *Internal) SetNetRpc(netRpc rpc.IRpc) {
	agent.netRpc = netRpc
}

func (agent *Internal) HandlerAccept(server net.IServer, session net.ISession) {
	agent.loop.Post(func() {
		agent.eventHandler.OnAccept(server, session)
	})
}

func (agent *Internal) HandlerSessionData(server net.IServer, session net.ISession, data []byte) {
	rpk := packet.NewReadPacket(data, binary.BigEndian, 0, len(data))
	agent.netRpc.HandleSessionPacket(session, rpk)
}

func (agent *Internal) HandlerSessionClose(server net.IServer, session net.ISession) {
	agent.loop.Post(func() {
		agent.eventHandler.OnSessionClose(server, session)
	})
}
