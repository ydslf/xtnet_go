package test

import (
	"fmt"
	"testing"
)

func TestSlice(t *testing.T) {
	s := make([]byte, 0, 5)
	fmt.Println(len(s))
	s = s[:5]
	fmt.Println(len(s))
	s[0] = 1
	fmt.Println(s)

	a := make([]byte, 5, 5)
	b := []byte{1, 2, 3}
	copy(a[1:], b)
	fmt.Println(a, b)
}
