package tcp

import "io"

type PkgProcessor interface {
	Read(session *Session) (data []byte, err error)
	Write(session *Session, data []byte) error
}

type PkgProcessorDefault struct {
}

func (process *PkgProcessorDefault) Read(session *Session) ([]byte, error) {
	buffLen := make([]byte, 4)
	if _, err := io.ReadFull(session.conn, buffLen); err != nil {
		return nil, err
	}
	return buffLen, nil
}
