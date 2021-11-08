package tcp

import (
	"net"
	xtNet "xtnet/net"
)

type ISession interface {
	xtNet.ISession
	setPktProc(IPktProc)
	start()
}

type ISessionCreator interface {
	CreateSession(xtNet.INetBase, net.Conn) ISession
}

type IPktProc interface {
	UnPack(session *Session) ([]byte, error)
	Pack(data []byte) []byte
}

type IPktProcCreator interface {
	CreatePktProc() IPktProc
}
