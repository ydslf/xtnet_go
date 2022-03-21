package util

import "errors"

var (
	BufErrOutOfRange = errors.New("buffer out of range")
)

type Buffer struct {
	data     []byte
	dataSize int
	writePos int
	readPos  int
}

func NewBuffer() *Buffer {
	return &Buffer{}
}

func (b *Buffer) GetWritePos() int {
	return b.writePos
}

func (b *Buffer) GetReadPos() int {
	return b.readPos
}

func (b *Buffer) GetWriteValidSize() int {
	return b.dataSize - b.writePos
}

func (b *Buffer) GetReadValidSize() int {
	return b.writePos - b.readPos
}

func (b *Buffer) GetWriteData() []byte {
	return b.data[b.writePos:]
}

func (b *Buffer) GetReadData() []byte {
	return b.data[b.readPos:b.writePos]
}

func (b *Buffer) AddWritePos(v int) error {
	p := b.writePos + v
	if p > b.dataSize {
		return BufErrOutOfRange
	}
	b.writePos = p
	return nil
}

func (b *Buffer) AddReadPos(v int) error {
	p := b.readPos + v
	if p > b.writePos {
		return BufErrOutOfRange
	}
	b.readPos = p
	return nil
}

func (b *Buffer) AdjustToHead() {
	if b.readPos > 0 {
		canReadSize := b.GetReadValidSize()
		if canReadSize > 0 {
			copy(b.data, b.data[b.readPos:b.writePos])
		}
		b.writePos = canReadSize
		b.readPos = 0
	}
}

func (b *Buffer) Write(data []byte) error {
	l := len(data)
	if b.GetWriteValidSize() < l {
		return BufErrOutOfRange
	}

	copy(b.data[b.writePos:], data)
	b.writePos += l

	return nil
}

func (b *Buffer) expand(size int) {
	size = int(SizeOfPow2(uint32(size)))
	if size < 64 {
		size = 64
	}

	dataSizeNew := b.dataSize + size
	dataNew := make([]byte, dataSizeNew, dataSizeNew)
	copy(dataNew, b.data[:b.writePos])
	b.data = dataNew
	b.dataSize = dataSizeNew
}

func (b *Buffer) MakeSureWriteEnough(size int) {
	if b.GetWriteValidSize() < size {
		b.AdjustToHead()
		if b.GetWriteValidSize() < size {
			b.expand(size - b.GetWriteValidSize())
		}
	}
}

func (b *Buffer) Reset() {
	b.writePos = 0
	b.readPos = 0
}
