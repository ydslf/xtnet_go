package session

import (
	"encoding/binary"
	"xtnet/frame"
	"xtnet/net"
	"xtnet/net/eventhandler"
	"xtnet/net/packet"
	"xtnet/net/rpc"
)

type Agent struct {
	loop         *frame.Loop
	netRpc       *rpc.Rpc
	eventHandler *eventhandler.Session
}

func NewAgent(service *frame.Service) *Agent {
	return &Agent{
		loop: service.GetLoop(),
	}
}

func (agent *Agent) SetNetRpc(netRpc *rpc.Rpc) {
	agent.netRpc = netRpc
}

func (agent *Agent) HandlerAccept(session net.ISession) {
	agent.loop.Post(func() {
		agent.eventHandler.OnAccept(session)
	})
}

func (agent *Agent) HandlerSessionClose(session net.ISession) {
	agent.loop.Post(func() {
		agent.eventHandler.OnSessionClose(session)
	})
}

func (agent *Agent) HandlerSessionData(session net.ISession, data []byte) {
	rpk := packet.NewReadPacket(data, binary.BigEndian, 0, len(data))
	if agent.netRpc != nil {
		agent.netRpc.HandleSessionPacket(session, rpk)
	} else {
		agent.loop.Post(func() {
			agent.eventHandler.OnSessionPacket(session, rpk)
		})
	}
}
