package tcp

import (
	"net"
	"sync"
)

type Session struct {
	conn			net.Conn
	onSessionData	OnSessionData
	onSessionClose	OnSessionClose
	close			bool
	wgClose			sync.WaitGroup
}

func newSession(conn net.Conn, onSessionData OnSessionData, onSessionClose OnSessionClose) *Session{
	return &Session{
		conn: conn,
		onSessionData: onSessionData,
		onSessionClose: onSessionClose,
		close: false,
	}
}

func (session *Session) Start(){
	go session.read()
	go session.write()
}

func (session *Session) Close(){

}

func (session *Session) read(){

}

func (session *Session) write(){

}