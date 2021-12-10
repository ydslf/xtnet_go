package net

type OnSessionStart func(ISession)
type OnSessionData func(ISession, []byte)
type OnSessionClose func(ISession)

type IServer interface {
	Start() bool
	Close()
}

type IClient interface {
	GetSession() ISession
	Connect() bool
	ConnectSync(ms int) error
	Send([]byte)
	Close(waitWrite bool)
}

type ISession interface {
	Send([]byte)
	Close(waitWrite bool)
	SetSessionStartCb(OnSessionStart)
	SetSessionDataCb(OnSessionData)
	SetSessionCloseCb(OnSessionClose)
	SetUserData(userData interface{})
	GetUserData() interface{}
}

type IServerAgent interface {
	HandleAccept(IServer, ISession)
	HandleSessionData(IServer, ISession, []byte)
	HandleSessionClose(IServer, ISession)
}

type IClientAgent interface {
	HandleConnect(IClient)
	HandleDisconnect(IClient)
	HandleClientData(IClient, []byte)
	HandleConnectBreak(IClient)
}
