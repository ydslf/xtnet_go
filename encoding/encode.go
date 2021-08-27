package encoding

import (
	"encoding/binary"
	"math"
	"reflect"
	"xtnet/util"
)

func Encode(obj interface{}) ([]byte, error) {
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	size, _, err := dataSize(v)
	if err != nil {
		return nil, err
	}

	buf := make([]byte, size)
	e := &encoder{order: binary.BigEndian, buf: buf}
	e.write(v)

	return buf, nil
}

func Encode2Buf(buf *util.Buffer, obj interface{}) error {
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	size, _, err := dataSize(v)
	if err != nil {
		return err
	}

	buf.MakeSureWriteEnough(size)
	e := &encoder{order: binary.BigEndian, buf: buf.GetWriteData()}
	e.write(v)
	err = buf.AddWritePos(size)

	return err
}

type encoder struct {
	order  binary.ByteOrder
	buf    []byte
	offset int
}

func (e *encoder) write(v reflect.Value) {
	switch v.Kind() {
	case reflect.Bool:
		e.bool(v.Bool())

	case reflect.Int8:
		e.int8(int8(v.Int()))
	case reflect.Int16:
		e.int16(int16(v.Int()))
	case reflect.Int32:
		e.int32(int32(v.Int()))
	case reflect.Int64:
		e.int64(v.Int())

	case reflect.Uint8:
		e.uint8(uint8(v.Uint()))
	case reflect.Uint16:
		e.uint16(uint16(v.Uint()))
	case reflect.Uint32:
		e.uint32(uint32(v.Uint()))
	case reflect.Uint64:
		e.uint64(v.Uint())

	case reflect.Float32:
		e.uint32(math.Float32bits(float32(v.Float())))
	case reflect.Float64:
		e.uint64(math.Float64bits(v.Float()))

	case reflect.Complex64:
		x := v.Complex()
		e.uint32(math.Float32bits(float32(real(x))))
		e.uint32(math.Float32bits(float32(imag(x))))
	case reflect.Complex128:
		x := v.Complex()
		e.uint64(math.Float64bits(real(x)))
		e.uint64(math.Float64bits(imag(x)))

	case reflect.String:
		e.writeString(v)

	case reflect.Array:
		e.writeArray(v, false)

	case reflect.Slice:
		e.writeArray(v, true)

	case reflect.Struct:
		e.writeStruct(v)
	}
}

func (e *encoder) bool(x bool) {
	if x {
		e.buf[e.offset] = 1
	} else {
		e.buf[e.offset] = 0
	}
	e.offset++
}

func (e *encoder) uint8(x uint8) {
	e.buf[e.offset] = x
	e.offset++
}

func (e *encoder) uint16(x uint16) {
	e.order.PutUint16(e.buf[e.offset:e.offset+2], x)
	e.offset += 2
}

func (e *encoder) uint32(x uint32) {
	e.order.PutUint32(e.buf[e.offset:e.offset+4], x)
	e.offset += 4
}

func (e *encoder) uint64(x uint64) {
	e.order.PutUint64(e.buf[e.offset:e.offset+8], x)
	e.offset += 8
}

func (e *encoder) int8(x int8) { e.uint8(uint8(x)) }

func (e *encoder) int16(x int16) { e.uint16(uint16(x)) }

func (e *encoder) int32(x int32) { e.uint32(uint32(x)) }

func (e *encoder) int64(x int64) { e.uint64(uint64(x)) }

func (e *encoder) writePrefix(l uint64) {
	for l >= 0x80 {
		e.buf[e.offset] = uint8(l&0x7f | 0x80)
		l >>= 7
		e.offset++
	}
	e.buf[e.offset] = uint8(l)
	e.offset++
}

func (e *encoder) writeString(v reflect.Value) {
	str := v.String()
	l := len(str)
	e.writePrefix(uint64(l))
	if l > 0 {
		copy(e.buf[e.offset:], str)
		e.offset += l
	}
}

func (e *encoder) writeArray(v reflect.Value, prefix bool) {
	if prefix {
		e.writePrefix(uint64(v.Len()))
	}

	if k := v.Type().Elem().Kind(); k == reflect.Uint8 {
		bytes := v.Bytes()
		l := len(bytes)
		if l > 0 {
			copy(e.buf[e.offset:], bytes)
			e.offset += l
		}
		return
	}

	l := v.Len()
	for i := 0; i < l; i++ {
		e.write(v.Index(i))
	}
}
