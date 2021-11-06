package agent

import (
	"encoding/binary"
	myNet "xtnet/net"
)

//TODO loop改为service，service再设置一个msgHandler, HandlerSessionData的时候，调用service.msgHandler.OnSessionPacket
type External struct {
	byteOrder binary.ByteOrder
}

func NewExternal(byteOrder binary.ByteOrder) myNet.IAdapter {
	return &External{
		byteOrder: byteOrder,
	}
}

func (agent *External) HandlerAccept(session myNet.ISession) {

}

func (agent *External) HandlerConnected(session myNet.ISession) {

}

func (agent *External) HandlerSessionClose(session myNet.ISession) {

}

func (agent *External) HandlerSessionData(session myNet.ISession, data []byte) {

}
