package util

import (
	list2 "container/list"
	"fmt"
	"testing"
)

func TestSizeOfPow2(t *testing.T) {
	var size uint32 = 123456
	t.Log(SizeOfPow2(size))
}

func BenchmarkSizeOfPow2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SizeOfPow2(123456)
	}
}

func TestList(t *testing.T) {
	l := list2.New()
	e4 := l.PushBack(4)
	e1 := l.PushFront(1)
	l.InsertBefore(3, e4)
	l.InsertAfter(2, e1)
	l.Remove(e4)
	for e := l.Front(); e != nil; e = e.Next() {
		fmt.Println(e.Value)
	}
}

func TestLog(t *testing.T) {

}
