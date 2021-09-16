package tcp

type SessionOuterClient struct {
	Session
}

func (session *SessionOuterClient) start() {
	session.receiveKey()
	session.doStart()
}

func (session *SessionOuterClient) receiveKey() {
	//TODO
}
