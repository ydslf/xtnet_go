package timer

import (
	"fmt"
	"testing"
	"time"
	"xtnet/frame"
)

func TestSystemTimer(t *testing.T) {
	serviceMain := frame.NewService()
	loop := serviceMain.GetLoop()
	manager := NewManager(loop)
	timer := manager.NewTimer(System)
	timer.Start(time.Second*1, false, func() {
		fmt.Println("startedEnd", time.Now())
		fmt.Println(timer)
	})
	loop.Run()
}
