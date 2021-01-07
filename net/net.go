package net

type OnAccept func(Session)
type OnSessionData func(Session, []byte)
type OnSessionClose func(Session)

type Session interface {
	Send(data []byte)
}
