package frame

import (
	mynet "xtnet/net"
)

//处理msg, net中的rpc （service 有自己的RPC）

type NetAgent struct {
	loop           *Loop
	onAccept       mynet.OnAccept
	onSessionData  mynet.OnSessionData
	onSessionClose mynet.OnSessionClose
}

func NewNetAgent(loop *Loop) *NetAgent {
	return &NetAgent{
		loop: loop,
	}
}

func (agent *NetAgent) SetCbOnAccept(onAccept mynet.OnAccept) {
	agent.onAccept = onAccept
}

func (agent *NetAgent) SetCbOnSessionData(onSessionData mynet.OnSessionData) {
	agent.onSessionData = onSessionData
}

func (agent *NetAgent) SetCbOnSessionClose(onSessionClose mynet.OnSessionClose) {
	agent.onSessionClose = onSessionClose
}

func (agent *NetAgent) OnAccept(session mynet.Session) {
	agent.loop.Post(func() {
		agent.onAccept(session)
	})
}

func (agent *NetAgent) OnSessionData(session mynet.Session, data []byte) {

}

func (agent *NetAgent) OnSessionClose(session mynet.Session) {
	agent.loop.Post(func() {
		agent.onSessionClose(session)
	})
}
