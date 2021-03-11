package util

import (
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
