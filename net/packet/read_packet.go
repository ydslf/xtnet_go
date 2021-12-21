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
	pos        int
	posReverse int
}

func NewReadPacket(data []byte, order binary.ByteOrder, pos int, posReverse int) *ReadPacket {
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

func (rpk *ReadPacket) GetLeftSize() int {
	return rpk.posReverse - rpk.pos
}

func (rpk *ReadPacket) CheckSize(size int) bool {
	return rpk.GetLeftSize() >= size
}

func (rpk *ReadPacket) PeakData(size int) []byte {
	if rpk.CheckSize(size) {
		return rpk.data[rpk.pos : rpk.pos+size]
	}
	return nil
}

func (rpk *ReadPacket) ReadData(size int) []byte {
	var ret []byte
	if rpk.CheckSize(size) {
		ret = rpk.data[rpk.pos : rpk.pos+size]
		rpk.pos += size
	}
	return ret
}

func (rpk *ReadPacket) ReadString() string {
	size := int(rpk.ReadInt16())
	if size > 0 {
		if rpk.CheckSize(size) {
			ret := string(rpk.data[rpk.pos : rpk.pos+size])
			rpk.pos += size
			return ret
		}
	}
	return ""
}

func (rpk *ReadPacket) PeakUint8() uint8 {
	var ret uint8 = 0
	if rpk.CheckSize(int(unsafe.Sizeof(ret))) {
		ret = rpk.data[rpk.pos]
	}
	return ret
}

func (rpk *ReadPacket) PeakInt8() int8 {
	return int8(rpk.PeakUint8())
}

func (rpk *ReadPacket) PeakUint16() uint16 {
	var ret uint16 = 0
	if rpk.CheckSize(int(unsafe.Sizeof(ret))) {
		ret = rpk.order.Uint16(rpk.data[rpk.pos:])
	}
	return ret
}

func (rpk *ReadPacket) PeakInt16() int16 {
	return int16(rpk.PeakUint16())
}

func (rpk *ReadPacket) PeakUint32() uint32 {
	var ret uint32 = 0
	if rpk.CheckSize(int(unsafe.Sizeof(ret))) {
		ret = rpk.order.Uint32(rpk.data[rpk.pos:])
	}
	return ret
}

func (rpk *ReadPacket) PeakInt32() int32 {
	return int32(rpk.PeakUint32())
}

func (rpk *ReadPacket) PeakUint64() uint64 {
	var ret uint64 = 0
	if rpk.CheckSize(int(unsafe.Sizeof(ret))) {
		ret = rpk.order.Uint64(rpk.data[rpk.pos:])
	}
	return ret
}

func (rpk *ReadPacket) PeakInt64() int64 {
	return int64(rpk.PeakUint64())
}

func (rpk *ReadPacket) ReadUint8() uint8 {
	var ret uint8 = 0
	size := int(unsafe.Sizeof(ret))
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
	size := int(unsafe.Sizeof(ret))
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
	size := int(unsafe.Sizeof(ret))
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
	size := int(unsafe.Sizeof(ret))
	if rpk.CheckSize(size) {
		ret = rpk.order.Uint64(rpk.data[rpk.pos:])
		rpk.pos += size
	}
	return ret
}

func (rpk *ReadPacket) ReadInt64() int64 {
	return int64(rpk.ReadUint64())
}

//---------------------------------------------------------------------

func (rpk *ReadPacket) PeakRevUint8() uint8 {
	var ret uint8 = 0
	if rpk.CheckSize(int(unsafe.Sizeof(ret))) {
		ret = rpk.data[rpk.posReverse-1]
	}
	return ret
}

func (rpk *ReadPacket) PeakRevInt8() int8 {
	return int8(rpk.PeakRevUint8())
}

func (rpk *ReadPacket) PeakRevUint16() uint16 {
	var ret uint16 = 0
	if rpk.CheckSize(int(unsafe.Sizeof(ret))) {
		ret = rpk.order.Uint16(rpk.data[rpk.posReverse-2:])
	}
	return ret
}

func (rpk *ReadPacket) PeakRevInt16() int16 {
	return int16(rpk.PeakRevUint16())
}

func (rpk *ReadPacket) PeakRevUint32() uint32 {
	var ret uint32 = 0
	if rpk.CheckSize(int(unsafe.Sizeof(ret))) {
		ret = rpk.order.Uint32(rpk.data[rpk.posReverse-4:])
	}
	return ret
}

func (rpk *ReadPacket) PeakRevInt32() int32 {
	return int32(rpk.PeakRevUint32())
}

func (rpk *ReadPacket) PeakRevUint64() uint64 {
	var ret uint64 = 0
	if rpk.CheckSize(int(unsafe.Sizeof(ret))) {
		ret = rpk.order.Uint64(rpk.data[rpk.posReverse-8:])
	}
	return ret
}

func (rpk *ReadPacket) PeakRevInt64() int64 {
	return int64(rpk.PeakRevUint64())
}

func (rpk *ReadPacket) ReadRevUint8() uint8 {
	var ret uint8 = 0
	size := int(unsafe.Sizeof(ret))
	if rpk.CheckSize(size) {
		ret = rpk.data[rpk.posReverse-size]
		rpk.pos += size
	}
	return ret
}

func (rpk *ReadPacket) ReadRevInt8() int8 {
	return int8(rpk.ReadRevUint8())
}

func (rpk *ReadPacket) ReadRevUint16() uint16 {
	var ret uint16 = 0
	size := int(unsafe.Sizeof(ret))
	if rpk.CheckSize(size) {
		ret = rpk.order.Uint16(rpk.data[rpk.posReverse-size:])
		rpk.pos += size
	}
	return ret
}

func (rpk *ReadPacket) ReadRevInt16() int16 {
	return int16(rpk.ReadRevUint16())
}

func (rpk *ReadPacket) ReadRevUint32() uint32 {
	var ret uint32 = 0
	size := int(unsafe.Sizeof(ret))
	if rpk.CheckSize(size) {
		ret = rpk.order.Uint32(rpk.data[rpk.posReverse-size:])
		rpk.pos += size
	}
	return ret
}

func (rpk *ReadPacket) ReadRevInt32() int32 {
	return int32(rpk.ReadRevUint32())
}

func (rpk *ReadPacket) ReadRevUint64() uint64 {
	var ret uint64 = 0
	size := int(unsafe.Sizeof(ret))
	if rpk.CheckSize(size) {
		ret = rpk.order.Uint64(rpk.data[rpk.posReverse-size:])
		rpk.pos += size
	}
	return ret
}

func (rpk *ReadPacket) ReadRevInt64() int64 {
	return int64(rpk.ReadRevUint64())
}
