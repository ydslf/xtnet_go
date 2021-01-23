package main

import (
	"encoding/binary"
	"fmt"
	"xtnet/net/packet"
)

func main() {
	data := []byte{1, 0, 0, 0}
	rpk := packet.NewReadPacket(data, binary.BigEndian, 0, 4)
	fmt.Println(rpk.PeakInt16())
}
