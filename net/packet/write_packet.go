package packet

import "encoding/binary"

type WritePacket struct {
	data  []byte
	order binary.ByteOrder
	pos   uint
}
