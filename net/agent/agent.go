package agent

import (
	"encoding/binary"
	"xtnet/frame"
	"xtnet/net"
	"xtnet/net/packet"
	"xtnet/net/rpc"
)

type Agent struct {
	loop            *frame.Loop
	netRpc          *rpc.Rpc
	onAccept        net.OnAccept
	onConnected     net.OnConnected
	onSessionPacket net.OnSessionPacket
	onSessionClose  net.OnSessionClose
}

func NewAgent(loop *frame.Loop) *Agent {
	return &Agent{
		loop: loop,
	}
}

func (agent *Agent) SetNetRpc(netRpc *rpc.Rpc) {
	agent.netRpc = netRpc
}

func (agent *Agent) SetOnAccept(onAccept net.OnAccept) {
	agent.onAccept = onAccept
}

func (agent *Agent) SetOnSessionPacket(onSessionPacket net.OnSessionPacket) {
	agent.onSessionPacket = onSessionPacket
}

func (agent *Agent) SetOnSessionClose(onSessionClose net.OnSessionClose) {
	agent.onSessionClose = onSessionClose
}

func (agent *Agent) HandlerAccept(session net.ISession) {
	agent.loop.Post(func() {
		agent.onAccept(session)
	})
}

func (agent *Agent) HandlerConnected(session net.ISession) {
	agent.loop.Post(func() {
		agent.onConnected(session)
	})
}

func (agent *Agent) HandlerSessionClose(session net.ISession) {
	agent.loop.Post(func() {
		agent.onSessionClose(session)
	})
}

func (agent *Agent) HandlerSessionData(session net.ISession, data []byte) {
	rpk := packet.NewReadPacket(data, binary.BigEndian, 0, len(data))
	if agent.netRpc != nil {
		agent.netRpc.HandleSessionPacket(session, rpk)
	} else {
		agent.loop.Post(func() {
			agent.onSessionPacket(session, rpk)
		})
	}
}
