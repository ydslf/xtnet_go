package tcp

import (
	"encoding/binary"
	"errors"
	"io"
)

type PktProcessor interface {
	Read(session *Session) (data []byte, err error)
	Write(session *Session, data []byte) error
}

var (
	NTErrPktTooLong = errors.New("packet too long")
	NTErrPktZero    = errors.New("packet len 0")
)

type PktProcessorDefault struct {
	maxPkgSize uint32
	byteOrder  binary.ByteOrder
}

/*
	| pkgHead | 	pkgBody 	|
	| pkgLen  | msgID |	msgBody |
*/

func NewPktProcessorDefault(maxPkgSize uint32, order binary.ByteOrder) PktProcessor {
	return &PktProcessorDefault{
		maxPkgSize: maxPkgSize,
		byteOrder:  order,
	}
}

func (process *PktProcessorDefault) Read(session *Session) ([]byte, error) {
	bufPktLen := make([]byte, 4)
	if _, err := io.ReadFull(session.conn, bufPktLen); err != nil {
		return nil, err
	}

	pktLen := process.byteOrder.Uint32(bufPktLen)
	if pktLen > process.maxPkgSize {
		//TODO
		return nil, NTErrPktTooLong
	}
	if pktLen == 0 {
		//TODO
		return nil, NTErrPktZero
	}

	pktData := make([]byte, pktLen)
	if _, err := io.ReadFull(session.conn, pktData); err != nil {
		return nil, err
	}

	return pktData, nil
}

func (process *PktProcessorDefault) Write(session *Session, data []byte) error {
	return nil
}
