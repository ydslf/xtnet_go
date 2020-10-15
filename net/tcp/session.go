package tcp

import (
	"fmt"
	"net"
	"sync"
)

type Session struct {
	conn         net.Conn
	pktProcessor PktProcessor
	close        bool
	wgClose      sync.WaitGroup
	sendChan     chan []byte

	onSessionData  OnSessionData
	onSessionClose OnSessionClose
}

func newSession(conn net.Conn, pktProcessor PktProcessor, sendBuffSize uint32) *Session {
	return &Session{
		conn:         conn,
		pktProcessor: pktProcessor,
		close:        false,
		sendChan:     make(chan []byte, sendBuffSize),
	}
}

func (session *Session) setCallback(onSessionData OnSessionData, onSessionClose OnSessionClose) {
	session.onSessionData = onSessionData
	session.onSessionClose = onSessionClose
}

func (session *Session) Send(data []byte) {
	if session.close {
		//TODO
		return
	}

	if len(session.sendChan) == cap(session.sendChan) {
		//TODO
		fmt.Printf("Session.sendChan is full")
	}
	session.sendChan <- data
}

func (session *Session) Close() {

}

func (session *Session) start() {
	go session.readRoutine()
	go session.writeRoutine()
}

func (session *Session) readRoutine() {
	for !session.close {
		data, err := session.pktProcessor.UnPack(session)
		if err != nil {
			//TODO
			return
		}

		session.onSessionData(session, data)
	}
}

func (session *Session) writeRoutine() {
	for data := range session.sendChan {
		if data == nil {
			break
		}

		pktData := session.pktProcessor.Pack(data)
		_, err := session.conn.Write(pktData)
		if err != nil {
			//TODO
			break
		}
	}
	//TODO
}
