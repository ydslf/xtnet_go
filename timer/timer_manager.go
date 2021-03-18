package timer

import (
	"time"
	"xtnet/frame"
)

const (
	TypeSystem int8 = 1 //系统timer
)

type Cb func()

type ITimer interface {
	Start(d time.Duration, repeat bool, cb Cb)
	Stop()
}

type Manager struct {
	loop *frame.Loop
}

func NewManager(loop *frame.Loop) *Manager {
	return &Manager{
		loop: loop,
	}
}

func (manager *Manager) NewTimer(timerType int8) ITimer {
	switch timerType {
	case TypeSystem:
		return NewSystemTimer(manager.loop)
	default:
		return nil
	}
}
