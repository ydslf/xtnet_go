package net

import "xtnet/frame"

type Rpc struct {
	loop *frame.Loop
}

func NewNetRpc(loop *frame.Loop) *Rpc {
	return &Rpc{
		loop: loop,
	}
}

func (netRpc *Rpc) HandleSessionData(session Session, data []byte) {

}

func (netRpc *Rpc) Request() {

}

func (netRpc *Rpc) Respond() {

}

func (netRpc *Rpc) Call() {

}

func (netRpc *Rpc) Answer() {

}
