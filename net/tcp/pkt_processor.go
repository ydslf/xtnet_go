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
	| pkgHead | 	pkgBody 	|
	| pkgLen  | msgID |	msgBody |
*/

const (
	maxPkgSizeDefault  uint32 = 1024 * 4
	pktHeadSizeDefault uint32 = 4
)

var orderDefault binary.ByteOrder = binary.BigEndian

type PktProcessor struct {
	pktHeadSize uint32
	maxPkgSize  uint32
	byteOrder   binary.ByteOrder
	headBuff    []byte
}

func (process *PktProcessor) UnPack(session *Session) ([]byte, error) {
	if _, err := io.ReadFull(session.conn, process.headBuff); err != nil {
		return nil, err
	}

	pktLen := process.byteOrder.Uint32(process.headBuff)
	if pktLen > process.maxPkgSize {
		xtnet_go.GetLogger().LogError("net.tcp.Unpack: pktLen > process.maxPkgSize, pktLen=%d, maxPkgSize=%d", pktLen, process.maxPkgSize)
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

func (process *PktProcessor) Pack(data []byte) []byte {
	//TODO 加一个写buffer，或者wpk直接传到session的写协程
	pktLen := uint32(len(data))
	pktData := make([]byte, process.pktHeadSize+pktLen)
	process.byteOrder.PutUint32(pktData, pktLen)
	copy(pktData[process.pktHeadSize:], data)
	return pktData
}

type PktProcessorMaker struct {
	pktHeadSize uint32
	maxPkgSize  uint32
	byteOrder   binary.ByteOrder
}

func NewPktProcessorMaker(pktHeadSize uint32, maxPkgSize uint32, order binary.ByteOrder) IPktProcessorMaker {
	return &PktProcessorMaker{
		pktHeadSize: pktHeadSize,
		maxPkgSize:  maxPkgSize,
		byteOrder:   order,
	}
}

func (m *PktProcessorMaker) CreatePktProcessor() IPktProcessor {
	return &PktProcessor{
		pktHeadSize: m.pktHeadSize,
		maxPkgSize:  m.maxPkgSize,
		byteOrder:   m.byteOrder,
		headBuff:    make([]byte, m.pktHeadSize, m.pktHeadSize),
	}
}
