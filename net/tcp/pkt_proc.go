package tcp

import (
	"encoding/binary"
	"errors"
	"io"
	xtnet_go "xtnet"
	"xtnet/util"
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
	maxPkgSizeDefault int = 1024 * 4
	pktHeadSize       int = 4
)

var orderDefault binary.ByteOrder = binary.BigEndian

type PktProc struct {
	maxPkgSize int
	byteOrder  binary.ByteOrder
	headBuff   []byte
	sendBuff   util.Buffer
}

func (proc *PktProc) UnPack(session *Session) ([]byte, error) {
	if _, err := io.ReadFull(session.conn, proc.headBuff); err != nil {
		return nil, err
	}

	pktLen := int(proc.byteOrder.Uint32(proc.headBuff))
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
	pktLen := len(data)
	proc.sendBuff.Reset()
	proc.sendBuff.MakeSureWriteEnough(pktHeadSize + pktLen)
	pktData := make([]byte, pktHeadSize+pktLen)
	proc.byteOrder.PutUint32(pktData, uint32(pktLen))
	copy(pktData[pktHeadSize:], data)
	return pktData
}

type PktProcCreator struct {
	maxPkgSize int
	byteOrder  binary.ByteOrder
}

func NewPktProcCreator(maxPkgSize int, order binary.ByteOrder) IPktProcCreator {
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
