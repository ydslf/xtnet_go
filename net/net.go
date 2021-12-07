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
}

type IServerAgent interface {
	HandlerAccept(IServer, ISession)
	HandlerSessionData(IServer, ISession, []byte)
	HandlerSessionClose(IServer, ISession)
}

type IClientAgent interface {
	HandlerConnect(IClient)
	HandlerDisconnect(IClient)
	HandlerClientData(IClient, []byte)
	HandlerConnectBreak(IClient)
}
