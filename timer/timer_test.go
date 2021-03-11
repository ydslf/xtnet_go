package timer

import (
	"fmt"
	"testing"
	"time"
)

func TestSystemTimer(t *testing.T) {
	AfterFunc(time.Millisecond*1, func() {
		fmt.Println("hahah")
	})

	time.Sleep(time.Hour * 1)
}
