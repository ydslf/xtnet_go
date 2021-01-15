package net

import "xtnet/frame"

type Agent struct {
	loop           *frame.Loop
	netRpc         *Rpc
	onAccept       OnAccept
	onSessionClose OnSessionClose
}

func NewNetAgent(loop *frame.Loop) *Agent {
	return &Agent{
		loop: loop,
	}
}

func (agent *Agent) SetNetRpc(netRpc *Rpc) {
	agent.netRpc = netRpc
}

func (agent *Agent) SetCbOnAccept(onAccept OnAccept) {
	agent.onAccept = onAccept
}

func (agent *Agent) SetCbOnSessionClose(onSessionClose OnSessionClose) {
	agent.onSessionClose = onSessionClose
}

func (agent *Agent) OnAccept(session Session) {
	agent.loop.Post(func() {
		agent.onAccept(session)
	})
}

func (agent *Agent) OnSessionClose(session Session) {
	agent.loop.Post(func() {
		agent.onSessionClose(session)
	})
}

func (agent *Agent) OnSessionData(session Session, data []byte) {
	if agent.netRpc != nil {

	} else {
		agent.loop.Post(func() {
			agent.HandleSessionData(session, data)
		})
	}
}

func (agent *Agent) HandleSessionData(session Session, data []byte) {

}
