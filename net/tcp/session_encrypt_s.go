package tcp

import (
	"net"
	"sync/atomic"
)

type SessionEncryptServer struct {
	Session
}

func (session *SessionEncryptServer) start() {
	//TODO 发送key, 开写协程写, 完了再开读协程
	if atomic.CompareAndSwapInt32(&session.status, sessionStatusNone, sessionStatusInit) {
		session.sendKey()
	}
}

func (session *SessionEncryptServer) sendKey() {
	//TODO
}

type SessionEncryptSCreator struct {
	sendChanSize int
}

func NewSessionEncryptSCreator(sendChanSize int) ISessionCreator {
	return &SessionEncryptSCreator{
		sendChanSize: sendChanSize,
	}
}

func (c *SessionEncryptSCreator) CreateSession(conn net.Conn) ISession {
	return &SessionEncryptServer{
		Session{
			conn:      conn,
			status:    sessionStatusNone,
			sendChan:  make(chan []byte, c.sendChanSize),
			closeChan: make(chan int, 1),
		},
	}
}
