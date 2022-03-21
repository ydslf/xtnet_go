package websocket

import (
	"github.com/gorilla/websocket"
	"sync"
	"sync/atomic"
	"time"
	xtnet "xtnet"
	xtnetNet "xtnet/net"
)

type closeType int

const (
	ctNone closeType = iota
	ctActive
	ctByRead
	ctByWrite
)

const (
	sendChanSizeDefault int = 1024
)

const (
	sessionStatusNone int32 = iota
	sessionStatusInit
	sessionStatusRunning
	sessionStatusClosing
	sessionStatusClosed
)

type Session struct {
	conn           *websocket.Conn
	status         int32
	wgClose        sync.WaitGroup
	sendChan       chan []byte
	closeChan      chan int
	closeType      closeType
	onSessionStart xtnetNet.OnSessionStart
	onSessionData  xtnetNet.OnSessionData
	onSessionClose xtnetNet.OnSessionClose
	userData       interface{}
}

func (session *Session) SetSessionStartCb(cb xtnetNet.OnSessionStart) {
	session.onSessionStart = cb
}

func (session *Session) SetSessionDataCb(cb xtnetNet.OnSessionData) {
	session.onSessionData = cb
}

func (session *Session) SetSessionCloseCb(cb xtnetNet.OnSessionClose) {
	session.onSessionClose = cb
}

func (session *Session) SetUserData(userData interface{}) {
	session.userData = userData
}

func (session *Session) GetUserData() interface{} {
	return session.userData
}

func (session *Session) GetRemoteAddr() string {
	return session.conn.RemoteAddr().String()
}

func (session *Session) Send(data []byte) {
	if atomic.LoadInt32(&session.status) != sessionStatusRunning {
		xtnet.GetLogger().LogWarn("websocket.Session.Send: session is not running")
		return
	}

	if len(session.sendChan) == cap(session.sendChan) {
		xtnet.GetLogger().LogWarn("websocket.Session.Send: sendChan is full")
	}

	session.sendChan <- data
}

//关闭session
func (session *Session) Close(waitWrite bool) {
	session.doClose(ctActive, waitWrite)
}

func (session *Session) doClose(ct closeType, waitWrite bool) bool {
	switch s := atomic.LoadInt32(&session.status); s {
	case sessionStatusInit, sessionStatusRunning:
		atomic.StoreInt32(&session.status, sessionStatusClosing)
		session.closeType = ct

		if ct == ctActive {
			session.conn.SetReadDeadline(time.Now())
			if waitWrite {
				close(session.sendChan)
			} else {
				session.closeChan <- 1
			}
		} else {
			if ct == ctByRead {
				session.closeChan <- 1
			} else {
				session.conn.SetReadDeadline(time.Now())
			}
		}
		return true
	default:
		if ct == ctActive {
			xtnet.GetLogger().LogWarn("websocket.Session.Close: session status error, status=%d", s)
		}
		return false
	}
}

func (session *Session) start() {
	if atomic.CompareAndSwapInt32(&session.status, sessionStatusNone, sessionStatusInit) {
		atomic.StoreInt32(&session.status, sessionStatusRunning)
		session.onSessionStart(session)

		session.wgClose.Add(2)
		go session.subRoutine()
		go session.readRoutine()
		go session.writeRoutine()
	}
}

func (session *Session) subRoutine() {
	session.wgClose.Wait()
	atomic.StoreInt32(&session.status, sessionStatusClosed)
	session.conn.Close()
	if session.closeType == ctByRead || session.closeType == ctByWrite {
		session.onSessionClose(session)
	}
}

func (session *Session) readRoutine() {
	defer session.wgClose.Done()

	for {
		_, data, err := session.conn.ReadMessage()
		if err != nil {
			//xtnet.GetLogger().LogError("websocket.session.readRoutine: err=%v", err)
			session.doClose(ctByRead, false)
			return
		}
		session.onSessionData(session, data)
	}
}

func (session *Session) writeRoutine() {
	defer session.wgClose.Done()

	for {
		select {
		case data, ok := <-session.sendChan:
			if ok == false {
				return
			}

			if err := session.conn.WriteMessage(websocket.BinaryMessage, data); err != nil {
				xtnet.GetLogger().LogError("websocket.session.writeRoutine: err=%v", err)
				session.doClose(ctByWrite, false)
				return
			}

		case <-session.closeChan:
			return
		}
	}
}

func NewSession(conn *websocket.Conn) *Session {
	return &Session{
		conn:      conn,
		status:    sessionStatusNone,
		sendChan:  make(chan []byte, sendChanSizeDefault),
		closeChan: make(chan int, 1),
		closeType: ctNone,
	}
}
