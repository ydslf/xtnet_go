package timer

import (
	"fmt"
	"testing"
	"time"
	"xtnet/frame"
)

func TestSystemTimer(t *testing.T) {
	loop := frame.NewLoop(0, true)
	manager := NewManager(loop)
	timer := manager.NewTimer(System)
	timer.Start(time.Second*1, 0, func() {
		fmt.Println("startedEnd", time.Now())
		fmt.Println(timer)
	})
	loop.RunOnce()
}

func TestSystemTimerRepeat(t *testing.T) {
	loop := frame.NewLoop(0, true)
	manager := NewManager(loop)
	timer := manager.NewTimer(System)
	timer.Start(time.Second*1, 2, func() {
		fmt.Println("startedEnd", time.Now())
		fmt.Println(timer)
	})
	for i := 0; i < 3; i++ {
		loop.RunOnce()
	}
}
