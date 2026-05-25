package timer

import (
	"time"
	"xtnet/frame"
)

const (
	System int8 = 1 //系统timer
)

const RepeatInfinity int = -1

type Cb func()

type ITimer interface {
	// Start repeat:
	// -1 = 无限循环
	//  0 = 重复0次,执行1次
	//  1 = 重复1次,执行2次
	//  2 = 重复2次,执行3次
	Start(d time.Duration, repeat int, cb Cb)
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
	case System:
		return NewSystemTimer(manager.loop)
	default:
		return nil
	}
}

func (manager *Manager) NewSystemTimer() ITimer {
	return NewSystemTimer(manager.loop)
}
