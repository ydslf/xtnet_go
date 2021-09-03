package main

import (
	"fmt"
	"time"
)

//channel 读取时堵塞操作，加了select也是
//关闭channel: 在channel的队列的前面的数据读完之后, 读取操作会立刻返回 零值和false (for range 也是)
//select 中，会忽略chan为nil的读取和写入

func channelTest1() {
	messages := make(chan string, 1)
	signals := make(chan bool)

	msg := "hi"
	select {
	case messages <- msg:
		fmt.Println("sent message", msg)
	default:
		fmt.Println("no message sent")
	}

	select {
	case msg, ok := <-messages:
		fmt.Println("received message", msg, ok)
	default:
		fmt.Println("no message received")
	}

	select {
	case msg, ok := <-messages:
		fmt.Println("received message", msg, ok)
	case sig := <-signals:
		fmt.Println("received signal", sig)
	default:
		fmt.Println("no activity")
	}
}

func channelTest2() {
	ch := make(chan string, 3)
	ch <- "a"
	ch <- "b"
	ch <- "c"
	close(ch)

	//ch <- "d"

	time.Sleep(time.Second * 1)
FOR:
	for {
		select {
		case value, ok := <-ch:
			if !ok {
				fmt.Println("quit", ok)
				break FOR
			}
			fmt.Println(value)
		}
	}

	for value := range ch {
		fmt.Println(value)
	}

	for i := 0; i < 5; i++ {
		value := <-ch
		fmt.Println("for1", value)
	}

	for i := 0; i < 5; i++ {
		value, ok := <-ch
		fmt.Println("for2", value, ok)
	}
}

//close channel
func channelTest3() {
	ch := make(chan string, 3)
	ch <- "a"
	ch <- "b"
	ch <- "c"

	go func() {
		for i := 0; i < 10; i++ {
			value, ok := <-ch
			fmt.Println(value, ok, time.Now())
		}
	}()

	time.Sleep(time.Second * 3)
	ch <- "d"
	time.Sleep(time.Second * 3)
	close(ch)

	time.Sleep(time.Hour * 1)
}

func channelClose() {
	ch := make(chan string, 3)
	ch <- "a"
	ch <- "b"
	ch <- "c"

	close(ch)
	time.Sleep(time.Second * 3)

	go func() {
		for v := range ch {
			fmt.Println(v, time.Now())
		}
		fmt.Println("quit", time.Now())
	}()

	time.Sleep(time.Hour * 1)
}

//cpu 100%, 加default
func channelTest4() {
	ch := make(chan string, 3)
	ch1 := make(chan int)
	ch <- "a"
	ch <- "b"
	ch <- "c"

	go func() {
		for {
			select {
			case value := <-ch:
				fmt.Println(value, time.Now())
			case value := <-ch1:
				fmt.Println(value, time.Now())
				//default:
			}
		}
	}()

	ch1 <- 1

	time.Sleep(time.Hour * 1)
}

//select no for
func channelTest5() {
	ch := make(chan string, 3)
	//ch <- "a"
	//ch <- "b"
	//ch <- "c"

	select {
	case ch <- "a":
		fmt.Println("send ok")
	}

	go func() {
		select {
		case value, ok := <-ch:
			fmt.Println(value, ok)
		}
	}()

	time.Sleep(time.Hour * 1)
}

func main() {
	//channelTest1()
	//channelTest2()
	//channelTest3()
	//channelTest4()
	//channelTest5()
	channelClose()
}
