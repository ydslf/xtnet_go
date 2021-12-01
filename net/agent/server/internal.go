package server

import (
	"xtnet/net"
)

type Internal struct {
}

func NewInternal() net.IServerAgent {
	return &Internal{}
}

func (agent *Internal) HandlerAccept(server net.IServer, session net.ISession) {

}

func (agent *Internal) HandlerSessionClose(server net.IServer, session net.ISession) {

}

func (agent *Internal) HandlerSessionData(server net.IServer, session net.ISession, data []byte) {

}
