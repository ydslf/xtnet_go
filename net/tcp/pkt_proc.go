package tcp

import (
	"encoding/binary"
	"errors"
	"io"
	xtnet_go "xtnet"
)

var (
	NTErrPktTooLong = errors.New("packet too long")
	NTErrPktZero    = errors.New("packet len 0")
)

//packet format
/*
	|           pktHead         |                      pktBody                     |
	| pktLen | crc32 | sequence |                      pktBody                     |
	|    pktLen | sequence      |                      pktBody                     |
*/

const (
	maxPkgSizeDefault uint32 = 1024 * 4
	pktHeadSize       uint32 = 4
)

var orderDefault binary.ByteOrder = binary.BigEndian

type PktProc struct {
	pktHeadSize uint32
	maxPkgSize  uint32
	byteOrder   binary.ByteOrder
	headBuff    []byte
}

func (proc *PktProc) UnPack(session *Session) ([]byte, error) {
	if _, err := io.ReadFull(session.conn, proc.headBuff); err != nil {
		return nil, err
	}

	pktLen := proc.byteOrder.Uint32(proc.headBuff)
	if pktLen > proc.maxPkgSize {
		xtnet_go.GetLogger().LogError("net.tcp.Unpack: pktLen > process.maxPkgSize, pktLen=%d, maxPkgSize=%d", pktLen, proc.maxPkgSize)
		return nil, NTErrPktTooLong
	}
	if pktLen == 0 {
		xtnet_go.GetLogger().LogError("net.tcp.Unpack: pktLen=0")
		return nil, NTErrPktZero
	}

	pktData := make([]byte, pktLen)
	if _, err := io.ReadFull(session.conn, pktData); err != nil {
		return nil, err
	}

	return pktData, nil
}

func (proc *PktProc) Pack(data []byte) []byte {
	//TODO 加一个写buffer，或者wpk直接传到session的写协程
	pktLen := uint32(len(data))
	pktData := make([]byte, proc.pktHeadSize+pktLen)
	proc.byteOrder.PutUint32(pktData, pktLen)
	copy(pktData[proc.pktHeadSize:], data)
	return pktData
}

type PktProcCreator struct {
	maxPkgSize uint32
	byteOrder  binary.ByteOrder
}

func NewPktProcCreator(maxPkgSize uint32, order binary.ByteOrder) IPktProcCreator {
	return &PktProcCreator{
		maxPkgSize: maxPkgSize,
		byteOrder:  order,
	}
}

func (c *PktProcCreator) CreatePktProc() IPktProc {
	return &PktProc{
		maxPkgSize: c.maxPkgSize,
		byteOrder:  c.byteOrder,
		headBuff:   make([]byte, pktHeadSize, pktHeadSize),
	}
}
