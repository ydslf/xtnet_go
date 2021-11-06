package agent

import (
	"encoding/binary"
	myNet "xtnet/net"
)

type Internal struct {
	byteOrder binary.ByteOrder
}

func NewInternal() myNet.IAdapter {
	return &Internal{
		binary.BigEndian,
	}
}

func (agent *Internal) HandlerAccept(session myNet.ISession) {

}

func (agent *Internal) HandlerConnected(session myNet.ISession) {

}

func (agent *Internal) HandlerSessionClose(session myNet.ISession) {

}

func (agent *Internal) HandlerSessionData(session myNet.ISession, data []byte) {

}
