package tcp

import (
	"encoding/binary"
	"errors"
	"io"
)

type PktProcessor interface {
	UnPack(session *Session) ([]byte, error)
	Pack(data []byte) []byte
}

var (
	NTErrPktTooLong = errors.New("packet too long")
	NTErrPktZero    = errors.New("packet len 0")
)

type PktProcessorDefault struct {
	pktHeadSize uint32
	maxPkgSize  uint32
	byteOrder   binary.ByteOrder
}

//PktProcessorDefault
/*
	| pkgHead | 	pkgBody 	|
	| pkgLen  | msgID |	msgBody |
*/

func NewPktProcessorDefault(maxPkgSize uint32, order binary.ByteOrder) PktProcessor {
	return &PktProcessorDefault{
		pktHeadSize: 4,
		maxPkgSize:  maxPkgSize,
		byteOrder:   order,
	}
}

func (process *PktProcessorDefault) UnPack(session *Session) ([]byte, error) {
	//TODO 优化 session中设置recieve buffer
	pktLenBuff := make([]byte, process.pktHeadSize)
	if _, err := io.ReadFull(session.conn, pktLenBuff); err != nil {
		return nil, err
	}

	pktLen := process.byteOrder.Uint32(pktLenBuff)
	if pktLen > process.maxPkgSize {
		//TODO
		//LOG
		return nil, NTErrPktTooLong
	}
	if pktLen == 0 {
		//TODO
		//LOG
		return nil, NTErrPktZero
	}

	pktData := make([]byte, pktLen)
	if _, err := io.ReadFull(session.conn, pktData); err != nil {
		return nil, err
	}

	return pktData, nil
}

func (process *PktProcessorDefault) Pack(data []byte) []byte {
	pktLen := uint32(len(data))
	pktData := make([]byte, process.pktHeadSize+pktLen)
	process.byteOrder.PutUint32(pktData, pktLen)
	copy(pktData[process.pktHeadSize:], data)
	return pktData
}
