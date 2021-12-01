package server

import (
	"encoding/binary"
	"xtnet/frame"
	"xtnet/net"
	"xtnet/net/eventhandler"
	"xtnet/net/packet"
	"xtnet/net/rpc"
)

type Normal struct {
	loop         *frame.Loop
	netRpc       *rpc.Rpc
	eventHandler *eventhandler.Server
}

func NewNormal(service *frame.Service) *Normal {
	return &Normal{
		loop: service.GetLoop(),
	}
}

func (agent *Normal) SetNetRpc(netRpc *rpc.Rpc) {
	agent.netRpc = netRpc
}

func (agent *Normal) HandlerAccept(server net.IServer, session net.ISession) {
	agent.loop.Post(func() {
		agent.eventHandler.OnAccept(server, session)
	})
}

func (agent *Normal) HandlerSessionData(server net.IServer, session net.ISession, data []byte) {
	rpk := packet.NewReadPacket(data, binary.BigEndian, 0, len(data))
	if agent.netRpc != nil {
		agent.netRpc.HandleSessionPacket(session, rpk)
	} else {
		agent.loop.Post(func() {
			agent.eventHandler.OnSessionPacket(server, session, rpk)
		})
	}
}

func (agent *Normal) HandlerSessionClose(server net.IServer, session net.ISession) {
	agent.loop.Post(func() {
		agent.eventHandler.OnSessionClose(server, session)
	})
}
