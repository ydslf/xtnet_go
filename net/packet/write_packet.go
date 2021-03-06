package packet

import (
	"encoding/binary"
	"unsafe"
	"xtnet/util"
)

/*
	|-----------dataReal-----------|
	|-reserveSize-|------data------|
	|-reserveSize-|----dataSize----|
*/

var wpkMaxLen int = 0xFFFF

type WritePacket struct {
	dataSize    int
	reserveSize int
	order       binary.ByteOrder
	dataReal    []byte
	data        []byte
	pos         int
	posRR       int
}

func SetWpkMaxLen(size int) {
	if size > wpkMaxLen {
		wpkMaxLen = size
	}
}

func NewWritePacket(size int, reserveSize int, order binary.ByteOrder) *WritePacket {
	wpk := &WritePacket{
		dataSize:    size,
		reserveSize: reserveSize,
		order:       order,
		pos:         0,
		posRR:       reserveSize,
	}
	wpk.dataReal = make([]byte, size+reserveSize, size+reserveSize)
	wpk.data = wpk.dataReal[reserveSize:]
	return wpk
}

func (wpk *WritePacket) GetValidSize() int {
	return wpk.dataSize - wpk.pos
}

func (wpk *WritePacket) GetReserveValidSize() int {
	return wpk.posRR
}

func (wpk *WritePacket) GetData() []byte {
	return wpk.data[wpk.pos:]
}

func (wpk *WritePacket) GetDataBySize(size int) []byte {
	return wpk.data[wpk.pos : wpk.pos+size]
}

func (wpk *WritePacket) GetRealData() []byte {
	return wpk.dataReal[wpk.posRR : wpk.reserveSize+wpk.pos]
}

func (wpk *WritePacket) CheckSize(size int) bool {
	return wpk.GetValidSize() >= size
}

func (wpk *WritePacket) AddPos(len int) {
	wpk.pos += len
}

func (wpk *WritePacket) expand(size int) bool {
	size = int(util.SizeOfPow2(uint32(size)))
	if size < 64 {
		size = 64
	}
	if (wpk.dataSize + size) > wpkMaxLen {
		return false
	}

	sizeRealNew := wpk.reserveSize + wpk.dataSize + size
	dataRealNew := make([]byte, sizeRealNew, sizeRealNew)
	copy(dataRealNew[wpk.posRR:], wpk.dataReal[wpk.posRR:wpk.reserveSize+wpk.pos])
	wpk.dataReal = dataRealNew
	wpk.dataSize += size
	wpk.data = wpk.dataReal[wpk.reserveSize:]

	return true
}

func (wpk *WritePacket) expandReserve(size int) bool {
	if size < 16 {
		size = 16
	}

	sizeRealNew := wpk.reserveSize + wpk.dataSize + size
	dataRealNew := make([]byte, sizeRealNew, sizeRealNew)
	copy(dataRealNew[wpk.posRR+size:], wpk.dataReal[wpk.posRR:wpk.reserveSize+wpk.pos])
	wpk.dataReal = dataRealNew
	wpk.reserveSize += size
	wpk.posRR += size
	wpk.data = wpk.dataReal[wpk.reserveSize:]

	return true
}

func (wpk *WritePacket) MakeSureEnough(size int) bool {
	if wpk.GetValidSize() < size {
		if !wpk.expand(size - wpk.GetValidSize()) {
			return false
		}
	}
	return true
}

func (wpk *WritePacket) MakeSureReserveEnough(size int) bool {
	if wpk.GetReserveValidSize() < size {
		if !wpk.expandReserve(size - wpk.GetValidSize()) {
			return false
		}
	}
	return true
}

func (wpk *WritePacket) WriteData(data []byte) {
	size := len(data)
	if wpk.MakeSureEnough(size) {
		copy(wpk.data[wpk.pos:], data)
		wpk.pos += size
	}
}

func (wpk *WritePacket) WriteString(data string) {
	size := len(data)
	if wpk.MakeSureEnough(size + 2) {
		wpk.WriteInt16(int16(size))
		copy(wpk.data[wpk.pos:], data)
		wpk.pos += size
	}
}

func (wpk *WritePacket) WriteReserveData(data []byte) {
	size := len(data)
	if wpk.MakeSureReserveEnough(size) {
		wpk.posRR -= size
		copy(wpk.dataReal[wpk.posRR:], data)
	}
}

func (wpk *WritePacket) WriteUint8(value uint8) {
	size := int(unsafe.Sizeof(value))
	if wpk.MakeSureEnough(size) {
		wpk.data[wpk.pos] = value
		wpk.pos += size
	}
}

func (wpk *WritePacket) WriteInt8(value int8) {
	wpk.WriteUint8(uint8(value))
}

func (wpk *WritePacket) WriteUint16(value uint16) {
	size := int(unsafe.Sizeof(value))
	if wpk.MakeSureEnough(size) {
		wpk.order.PutUint16(wpk.data[wpk.pos:], value)
		wpk.pos += size
	}
}

func (wpk *WritePacket) WriteInt16(value int16) {
	wpk.WriteUint16(uint16(value))
}

func (wpk *WritePacket) WriteUint32(value uint32) {
	size := int(unsafe.Sizeof(value))
	if wpk.MakeSureEnough(size) {
		wpk.order.PutUint32(wpk.data[wpk.pos:], value)
		wpk.pos += size
	}
}

func (wpk *WritePacket) WriteInt32(value int32) {
	wpk.WriteUint32(uint32(value))
}

func (wpk *WritePacket) WriteUint64(value uint64) {
	size := int(unsafe.Sizeof(value))
	if wpk.MakeSureEnough(size) {
		wpk.order.PutUint64(wpk.data[wpk.pos:], value)
		wpk.pos += size
	}
}

func (wpk *WritePacket) WriteInt64(value int64) {
	wpk.WriteUint64(uint64(value))
}

func (wpk *WritePacket) WriteReserveUint8(value uint8) {
	size := int(unsafe.Sizeof(value))
	if wpk.MakeSureReserveEnough(size) {
		wpk.posRR -= size
		wpk.dataReal[wpk.posRR] = value
	}
}

func (wpk *WritePacket) WriteReserveInt8(value int8) {
	wpk.WriteReserveUint8(uint8(value))
}

func (wpk *WritePacket) WriteReserveUint16(value uint16) {
	size := int(unsafe.Sizeof(value))
	if wpk.MakeSureReserveEnough(size) {
		wpk.posRR -= size
		wpk.order.PutUint16(wpk.dataReal[wpk.posRR:], value)
	}
}

func (wpk *WritePacket) WriteReserveInt16(value int16) {
	wpk.WriteReserveUint16(uint16(value))
}

func (wpk *WritePacket) WriteReserveUint32(value uint32) {
	size := int(unsafe.Sizeof(value))
	if wpk.MakeSureReserveEnough(size) {
		wpk.posRR -= size
		wpk.order.PutUint32(wpk.dataReal[wpk.posRR:], value)
	}
}

func (wpk *WritePacket) WriteReserveInt32(value int32) {
	wpk.WriteReserveUint32(uint32(value))
}

func (wpk *WritePacket) WriteReserveUint64(value uint64) {
	size := int(unsafe.Sizeof(value))
	if wpk.MakeSureReserveEnough(size) {
		wpk.posRR -= size
		wpk.order.PutUint64(wpk.dataReal[wpk.posRR:], value)
	}
}

func (wpk *WritePacket) WriteReserveInt64(value int64) {
	wpk.WriteReserveUint64(uint64(value))
}
