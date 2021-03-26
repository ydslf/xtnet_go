package main

import (
	"fmt"
	"net"
	"strings"
	"sync"
)

func channelTest1() {
	messages := make(chan string)
	signals := make(chan bool)

	msg := "hi"
	select {
	case messages <- msg:
		fmt.Println("sent message", msg)
	default:
		fmt.Println("no message sent")
	}

	select {
	case msg := <-messages:
		fmt.Println("received message", msg)
	default:
		fmt.Println("no message received")
	}

	select {
	case msg := <-messages:
		fmt.Println("received message", msg)
	case sig := <-signals:
		fmt.Println("received signal", sig)
	default:
		fmt.Println("no activity")
	}
}

//--------------------------------------------------------------------------------------------
type TcpListeners struct {
	addrs     []string
	conns     chan *net.TCPConn
	closing   chan bool
	wait      *sync.WaitGroup
	listeners []*net.TCPListener
}

func NewTcpListeners(addrs []string) (v *TcpListeners, err error) {
	v = &TcpListeners{
		addrs:   addrs,
		conns:   make(chan *net.TCPConn),
		closing: make(chan bool, 1),
		wait:    &sync.WaitGroup{},
	}

	return
}

// Listen at addrs format as netowrk://laddr, for example,
// tcp://:1935, tcp4://:1935, tcp6://1935, tcp://0.0.0.0:1935
func (v *TcpListeners) ListenTCP() (err error) {
	for _, addr := range v.addrs {
		vs := strings.Split(addr, "://")
		network, laddr := vs[0], vs[1]

		if l, err := net.Listen(network, laddr); err != nil {
			return err
		} else {
			v.listeners = append(v.listeners, l.(*net.TCPListener))
		}
	}

	v.wait.Add(len(v.listeners))
	for i, l := range v.listeners {
		addr := v.addrs[i]
		go func(l *net.TCPListener, addr string) {
			defer v.wait.Done()
			for {
				var conn *net.TCPConn
				if conn, err = l.AcceptTCP(); err != nil {
					return
				}

				select {
				case v.conns <- conn:
				case c := <-v.closing:
					v.closing <- c
					conn.Close()
				}
			}
		}(l, addr)
	}

	return
}

func (v *TcpListeners) Close() (err error) {
	// unblock all listener internal goroutines
	v.closing <- true

	// interrupt all listeners.
	for _, v := range v.listeners {
		if r := v.Close(); r != nil {
			err = r
		}
	}

	// wait for all listener internal goroutines to quit.
	v.wait.Wait()

	// clear the closing signal.
	_ = <-v.closing

	// close channels to unblock the user goroutine to AcceptTCP()
	close(v.conns)

	return
}

func channelTest2() {
	ch := make(chan int, 3)
	ch <- 3
	ch <- 2
	ch <- 1
	close(ch)
	for i := 0; i < 5; i++ {
		value, ok := <-ch
		fmt.Println(value, ok)
	}
}

//输出：32100

func main() {
	//channelTest1()
	channelTest2()
}
