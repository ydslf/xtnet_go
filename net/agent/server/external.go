package server

import (
	"xtnet/net"
)

type External struct {
}

func NewExternal() net.IServerAgent {
	return &External{}
}

func (agent *External) HandlerAccept(server net.IServer, session net.ISession) {
}

func (agent *External) HandlerSessionClose(server net.IServer, session net.ISession) {
}

func (agent *External) HandlerSessionData(server net.IServer, session net.ISession, data []byte) {
	//rpk := packet.NewReadPacket(data, agent.byteOrder, 0, len(data))
	//serviceType := rpk.ReadInt8()
	//serviceID := rpk.ReadInt16()
}
