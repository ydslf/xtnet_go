package main

import (
	"fmt"
	"time"
	"xtnet/frame"
	xttimer "xtnet/timer"
)

func Test1() {
	for i := 0; i < 100; i++ {
		time.AfterFunc(time.Millisecond*1, func() {
			fmt.Println("timer1")
		})
	}
	time.Sleep(time.Hour * 1)
}

func Test2() {
	startedAt := time.Now()
	fmt.Println("startedAt", startedAt)
	timer1 := time.AfterFunc(time.Second*1, func() {
		fmt.Println("startedEnd", time.Now())
		fmt.Println("time pass", time.Since(startedAt))
	})
	//timer1.Stop()
	time.Sleep(time.Second * 2)
	fmt.Println("reset result:", timer1.Reset(time.Second*2))
	//timer1.C
	time.Sleep(time.Hour * 1)
}

func DoTicker(c <-chan time.Time) {
	startedAt := time.Now()
	fmt.Println("startedAt", startedAt)
	for i := range c {
		// do something
		fmt.Println("ticker return", i)
		fmt.Println("ticker pass", time.Since(startedAt))
	}
	fmt.Println("ticker stop")
}
func TestTicker() {
	ticker1 := time.NewTicker(time.Second * 1)
	go DoTicker(ticker1.C)
	time.Sleep(time.Second * 5)
	ticker1.Stop()
	time.Sleep(time.Hour * 1)
}

func TestSystemTimer() {
	serviceMain := frame.NewService()
	loop := serviceMain.GetLoop()
	manager := xttimer.NewManager(loop)
	timer := manager.NewTimer(xttimer.System)
	timer.Start(time.Second*5, false, func() {
		fmt.Println("startedEnd", time.Now())
		fmt.Println(timer)
	})
	timer1 := manager.NewTimer(xttimer.System)
	timer1.Start(time.Second*1, true, func() {
		fmt.Println("tick", time.Now())
		fmt.Println(timer1)
	})
	loop.Run()
}

func main() {
	//Test1()
	//Test2()
	//TestTicker()
	TestSystemTimer()
}
