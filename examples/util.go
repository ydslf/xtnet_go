package main

import (
	"fmt"
	"xtnet/util"
)

func main() {
	var size uint32 = 123456
	fmt.Println(util.SizeOfPow2(size))
}
