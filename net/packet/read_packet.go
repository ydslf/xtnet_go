package packet

import (
	"encoding/binary"
	"unsafe"
)

/*
	|----------data----------|
	[pos					 posReverse)
*/

type ReadPacket struct {
	data       []byte
	order      binary.ByteOrder
	pos        uint
	posReverse uint
}

func NewReadPacket(data []byte, order binary.ByteOrder, pos uint, posReverse uint) *ReadPacket {
	return &ReadPacket{
		data:       data,
		order:      order,
		pos:        pos,
		posReverse: posReverse,
	}
}

func (rpk *ReadPacket) GetCurData() []byte {
	return rpk.data[rpk.pos:rpk.posReverse]
}

func (rpk *ReadPacket) GetLeftSize() uint {
	return rpk.posReverse - rpk.pos
}

func (rpk *ReadPacket) CheckSize(size uint) bool {
	return rpk.GetLeftSize() >= size
}

func (rpk *ReadPacket) PeakData(size uint) []byte {
	if rpk.CheckSize(size) {
		return rpk.data[rpk.pos : rpk.pos+size]
	}
	return nil
}

func (rpk *ReadPacket) ReadData(size uint) []byte {
	var ret []byte
	if rpk.CheckSize(size) {
		ret = rpk.data[rpk.pos : rpk.pos+size]
		rpk.pos += size
	}
	return ret
}

func (rpk *ReadPacket) PeakUint8() uint8 {
	var ret uint8 = 0
	if rpk.CheckSize(uint(unsafe.Sizeof(ret))) {
		ret = rpk.data[rpk.pos]
	}
	return ret
}

func (rpk *ReadPacket) PeakInt8() int8 {
	return int8(rpk.PeakUint8())
}

func (rpk *ReadPacket) PeakUint16() uint16 {
	var ret uint16 = 0
	if rpk.CheckSize(uint(unsafe.Sizeof(ret))) {
		ret = rpk.order.Uint16(rpk.data[rpk.pos:])
	}
	return ret
}

func (rpk *ReadPacket) PeakInt16() int16 {
	return int16(rpk.PeakUint16())
}

func (rpk *ReadPacket) PeakUint32() uint32 {
	var ret uint32 = 0
	if rpk.CheckSize(uint(unsafe.Sizeof(ret))) {
		ret = rpk.order.Uint32(rpk.data[rpk.pos:])
	}
	return ret
}

func (rpk *ReadPacket) PeakInt32() int32 {
	return int32(rpk.PeakUint32())
}

func (rpk *ReadPacket) PeakUint64() uint64 {
	var ret uint64 = 0
	if rpk.CheckSize(uint(unsafe.Sizeof(ret))) {
		ret = rpk.order.Uint64(rpk.data[rpk.pos:])
	}
	return ret
}

func (rpk *ReadPacket) PeakInt64() int64 {
	return int64(rpk.PeakUint64())
}

func (rpk *ReadPacket) ReadUint8() uint8 {
	var ret uint8 = 0
	size := uint(unsafe.Sizeof(ret))
	if rpk.CheckSize(size) {
		ret = rpk.data[rpk.pos]
		rpk.pos += size
	}
	return ret
}

func (rpk *ReadPacket) ReadInt8() int8 {
	return int8(rpk.ReadUint8())
}

func (rpk *ReadPacket) ReadUint16() uint16 {
	var ret uint16 = 0
	size := uint(unsafe.Sizeof(ret))
	if rpk.CheckSize(size) {
		ret = rpk.order.Uint16(rpk.data[rpk.pos:])
		rpk.pos += size
	}
	return ret
}

func (rpk *ReadPacket) ReadInt16() int16 {
	return int16(rpk.ReadUint16())
}

func (rpk *ReadPacket) ReadUint32() uint32 {
	var ret uint32 = 0
	size := uint(unsafe.Sizeof(ret))
	if rpk.CheckSize(size) {
		ret = rpk.order.Uint32(rpk.data[rpk.pos:])
		rpk.pos += size
	}
	return ret
}

func (rpk *ReadPacket) ReadInt32() int32 {
	return int32(rpk.ReadUint32())
}

func (rpk *ReadPacket) ReadUint64() uint64 {
	var ret uint64 = 0
	size := uint(unsafe.Sizeof(ret))
	if rpk.CheckSize(size) {
		ret = rpk.order.Uint64(rpk.data[rpk.pos:])
		rpk.pos += size
	}
	return ret
}

func (rpk *ReadPacket) ReadInt64() int64 {
	return int64(rpk.ReadUint64())
}
