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

func TestWheelTimer(t *testing.T) {
	loop := frame.NewLoop(0, true)
	manager := NewManager(loop)
	timer := manager.NewTimer(Wheel)

	count := 0
	timer.Start(time.Millisecond*20, 0, func() {
		count++
	})
	loop.RunOnce()

	if count != 1 {
		t.Fatalf("wheel timer callback count=%d, want 1", count)
	}
}

func TestWheelTimerRepeat(t *testing.T) {
	loop := frame.NewLoop(0, true)
	manager := NewManager(loop)
	timer := manager.NewTimer(Wheel)

	count := 0
	timer.Start(time.Millisecond*20, 2, func() {
		count++
	})
	for i := 0; i < 3; i++ {
		loop.RunOnce()
	}

	if count != 3 {
		t.Fatalf("wheel timer callback count=%d, want 3", count)
	}
}

func TestWheelTimerStop(t *testing.T) {
	loop := frame.NewLoop(0, true)
	manager := NewManager(loop)
	timer := manager.NewTimer(Wheel)
	done := make(chan struct{}, 1)

	timer.Start(time.Millisecond*20, 0, func() {
		done <- struct{}{}
	})
	timer.Stop()
	go loop.Run()
	time.Sleep(time.Millisecond * 60)
	loop.Close(false)

	select {
	case <-done:
		t.Fatal("stopped wheel timer should not run callback")
	default:
	}
}

func TestWheelTimerRound(t *testing.T) {
	loop := frame.NewLoop(0, true)
	wheel := NewTimeWheel(loop, time.Millisecond*10, 2, 2)
	timer := wheel.NewTimer()

	// 70ms => ticks=7, totalTicks[last]=4 => round=1, remaining=3
	done := make(chan struct{}, 1)
	timer.Start(time.Millisecond*70, 0, func() {
		done <- struct{}{}
	})

	go loop.Run()
	defer loop.Close(false)

	select {
	case <-done:
	case <-time.After(time.Millisecond * 200):
		t.Fatal("round-based wheel timer did not fire")
	}
}

func TestWheelTimerRound1(t *testing.T) {
	loop := frame.NewLoop(0, true)
	wheel := NewTimeWheel(loop, time.Millisecond*10, 100, 60)
	timer := wheel.NewTimer()

	done := make(chan struct{}, 1)
	timer.Start(time.Millisecond*60001, 0, func() {
		done <- struct{}{}
	})

	go loop.Run()
	defer loop.Close(false)

	select {
	case <-done:
	case <-time.After(time.Second * 200):
		t.Fatal("round-based wheel timer did not fire")
	}
}

func TestWheelTimerZeroDuration(t *testing.T) {
	loop := frame.NewLoop(0, true)
	manager := NewManager(loop)
	timer := manager.NewTimer(Wheel)

	done := make(chan struct{}, 1)
	timer.Start(0, 0, func() {
		done <- struct{}{}
	})

	go loop.Run()
	defer loop.Close(false)

	select {
	case <-done:
	case <-time.After(time.Millisecond * 200):
		t.Fatal("zero-duration wheel timer did not fire")
	}
}
