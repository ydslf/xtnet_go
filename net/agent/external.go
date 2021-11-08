package agent

import (
	"encoding/binary"
	xtNet "xtnet/net"
)

//TODO loop改为service，service再设置一个msgHandler, HandlerSessionData的时候，调用service.msgHandler.OnSessionPacket
type External struct {
	byteOrder binary.ByteOrder
}

func NewExternal(byteOrder binary.ByteOrder) xtNet.IAgent {
	return &External{
		byteOrder: byteOrder,
	}
}

func (agent *External) HandlerAccept(session xtNet.ISession) {
}

func (agent *External) HandlerConnected(session xtNet.ISession) {
}

func (agent *External) HandlerSessionClose(session xtNet.ISession) {
}

func (agent *External) HandlerSessionData(session xtNet.ISession, data []byte) {
	//rpk := packet.NewReadPacket(data, agent.byteOrder, 0, len(data))
	//serviceType := rpk.ReadInt8()
	//serviceID := rpk.ReadInt16()
}
