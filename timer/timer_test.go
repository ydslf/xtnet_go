package timer

import (
	"fmt"
	"testing"
	"time"
	"xtnet/frame"
)

func TestSystemTimer(t *testing.T) {
	serviceMain := frame.NewService()
	manager := NewManager(serviceMain)
	timer := manager.NewTimer(System)
	timer.Start(time.Second*1, false, func() {
		fmt.Println("startedEnd", time.Now())
		fmt.Println(timer)
	})
	serviceMain.GetLoop().RunOnce()
}
