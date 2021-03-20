package packet

import (
	"encoding/binary"
	"fmt"
	"testing"
)

func TestReadPacket(t *testing.T) {
	data := []byte{1, 0, 0, 0}
	rpk := NewReadPacket(data, binary.BigEndian, 0, 4)
	fmt.Println(rpk.PeakInt16())
}
