package tcp

import (
	"net"
	myNet "xtnet/net"
)

type ISession interface {
	myNet.ISession
	setPktProc(IPktProc)
	start()
}

type ISessionCreator interface {
	CreateSession(myNet.INetBase, net.Conn, int) ISession
}

type IPktProc interface {
	UnPack(session *Session) ([]byte, error)
	Pack(data []byte) []byte
}

type IPktProcCreator interface {
	CreatePktProc() IPktProc
}
