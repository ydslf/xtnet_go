package tcp

type SessionOuterServer struct {
	Session
}

func (session *SessionOuterServer) start() {
	session.sendKey()
	session.doStart()
}

func (session *SessionOuterServer) sendKey() {
	//TODO
}
