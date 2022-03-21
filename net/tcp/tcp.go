package tcp

import (
	"net"
	xtnetNet "xtnet/net"
)

type ISession interface {
	xtnetNet.ISession
	setPktProc(IPktProc)
	start()
	doClose(closeType closeType, waitWrite bool) bool
}

type ISessionCreator interface {
	CreateSession(net.Conn) ISession
}

type IPktProc interface {
	UnPack(session *Session) ([]byte, error)
	Pack(data []byte) []byte
}

type IPktProcCreator interface {
	CreatePktProc() IPktProc
}
