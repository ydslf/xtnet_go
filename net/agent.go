package net

import (
	"encoding/binary"
	"xtnet/frame"
	"xtnet/net/packet"
)

type Agent struct {
	loop            *frame.Loop
	netRpc          *Rpc
	onAccept        OnAccept
	onSessionPacket OnSessionPacket
	onSessionClose  OnSessionClose
}

func NewAgent(loop *frame.Loop) *Agent {
	return &Agent{
		loop: loop,
	}
}

func (agent *Agent) SetNetRpc(netRpc *Rpc) {
	agent.netRpc = netRpc
}

func (agent *Agent) SetOnAccept(onAccept OnAccept) {
	agent.onAccept = onAccept
}

func (agent *Agent) SetOnSessionPacket(onSessionPacket OnSessionPacket) {
	agent.onSessionPacket = onSessionPacket
}

func (agent *Agent) SetOnSessionClose(onSessionClose OnSessionClose) {
	agent.onSessionClose = onSessionClose
}

func (agent *Agent) HandlerAccept(session ISession) {
	agent.loop.Post(func() {
		agent.onAccept(session)
	})
}

func (agent *Agent) HandlerConnected(session ISession) {

}

func (agent *Agent) HandlerSessionClose(session ISession) {
	agent.loop.Post(func() {
		agent.onSessionClose(session)
	})
}

func (agent *Agent) HandlerSessionData(session ISession, data []byte) {
	rpk := packet.NewReadPacket(data, binary.BigEndian, 0, uint(len(data)))
	if agent.netRpc != nil {
		agent.netRpc.HandleSessionPacket(session, rpk)
	} else {
		agent.loop.Post(func() {
			agent.onSessionPacket(session, rpk)
		})
	}
}
