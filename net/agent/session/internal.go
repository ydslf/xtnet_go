package session

import (
	"encoding/binary"
	xtNet "xtnet/net"
)

type Internal struct {
	byteOrder binary.ByteOrder
}

func NewInternal() xtNet.ISessionAgent {
	return &Internal{
		binary.BigEndian,
	}
}

func (agent *Internal) HandlerAccept(session xtNet.ISession) {

}

func (agent *Internal) HandlerConnected(session xtNet.ISession) {

}

func (agent *Internal) HandlerSessionClose(session xtNet.ISession) {

}

func (agent *Internal) HandlerSessionData(session xtNet.ISession, data []byte) {

}
