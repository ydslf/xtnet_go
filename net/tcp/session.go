package tcp

import (
	"fmt"
	"net"
	"sync"
	mynet "xtnet/net"
)

type Session struct {
	conn           net.Conn
	pktProcessor   IPktProcessor
	close          bool
	wgClose        sync.WaitGroup
	sendChan       chan []byte
	userData       interface{}
	onSessionData  mynet.OnSessionData
	onSessionClose mynet.OnSessionClose
	agent          mynet.IAgent
}

func newSession(conn net.Conn, pktProcessor IPktProcessor, sendBuffSize int) *Session {
	session := &Session{
		conn:         conn,
		pktProcessor: pktProcessor,
		close:        false,
		sendChan:     make(chan []byte, sendBuffSize),
	}

	return session
}

func (session *Session) setCallback(onSessionData mynet.OnSessionData, onSessionClose mynet.OnSessionClose) {
	session.onSessionData = onSessionData
	session.onSessionClose = onSessionClose
}

func (session *Session) setAgent(agent mynet.IAgent) {
	session.agent = agent
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

func (session *Session) SetUserData(userData interface{}) {
	session.userData = userData
}

func (session *Session) GetUserData() interface{} {
	return session.userData
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
			//LOG
			session.Close()
			return
		}

		if session.agent != nil {
			session.agent.HandlerSessionData(session, data)
		} else {
			session.onSessionData(session, data)
		}

	}
}

func (session *Session) writeRoutine() {
	for data := range session.sendChan {
		pktData := session.pktProcessor.Pack(data)
		_, err := session.conn.Write(pktData)
		if err != nil {
			//TODO
			//LOG
			//CLOSE
			break
		}
	}
	//TODO
}
