package encoding

import (
	"encoding/binary"
	"errors"
	"math"
	"reflect"
)

var (
	ECErrMustPtrType = errors.New("decode interface must be ptr")
)

func Decode(data []byte, obj interface{}) error {
	v := reflect.ValueOf(obj)

	if v.Kind() != reflect.Ptr {
		return ECErrMustPtrType
	}

	v = v.Elem()
	d := decoder{order: binary.BigEndian, buf: data}
	d.read(v)

	return nil
}

type decoder struct {
	order  binary.ByteOrder
	buf    []byte
	offset int
}

func (d *decoder) read(v reflect.Value) {
	switch v.Kind() {
	case reflect.Bool:
		v.SetBool(d.bool())

	case reflect.Int8:
		v.SetInt(int64(d.int8()))
	case reflect.Int16:
		v.SetInt(int64(d.int16()))
	case reflect.Int32:
		v.SetInt(int64(d.int32()))
	case reflect.Int64:
		v.SetInt(d.int64())

	case reflect.Uint8:
		v.SetUint(uint64(d.uint8()))
	case reflect.Uint16:
		v.SetUint(uint64(d.uint16()))
	case reflect.Uint32:
		v.SetUint(uint64(d.uint32()))
	case reflect.Uint64:
		v.SetUint(d.uint64())

	case reflect.Float32:
		v.SetFloat(float64(math.Float32frombits(d.uint32())))
	case reflect.Float64:
		v.SetFloat(math.Float64frombits(d.uint64()))

	case reflect.Complex64:
		v.SetComplex(complex(
			float64(math.Float32frombits(d.uint32())),
			float64(math.Float32frombits(d.uint32())),
		))
	case reflect.Complex128:
		v.SetComplex(complex(
			math.Float64frombits(d.uint64()),
			math.Float64frombits(d.uint64()),
		))
	case reflect.String:
		d.readString(v)

	case reflect.Array:
		d.readArray(v)

	case reflect.Slice:
		d.readSlice(v)

	case reflect.Struct:
		d.readStruct(v)
	}
}

func (d *decoder) bool() bool {
	x := d.buf[d.offset]
	d.offset++
	return x != 0
}

func (d *decoder) uint8() uint8 {
	x := d.buf[d.offset]
	d.offset++
	return x
}

func (d *decoder) uint16() uint16 {
	x := d.order.Uint16(d.buf[d.offset : d.offset+2])
	d.offset += 2
	return x
}

func (d *decoder) uint32() uint32 {
	x := d.order.Uint32(d.buf[d.offset : d.offset+4])
	d.offset += 4
	return x
}

func (d *decoder) uint64() uint64 {
	x := d.order.Uint64(d.buf[d.offset : d.offset+8])
	d.offset += 8
	return x
}

func (d *decoder) int8() int8 { return int8(d.uint8()) }

func (d *decoder) int16() int16 { return int16(d.uint16()) }

func (d *decoder) int32() int32 { return int32(d.uint32()) }

func (d *decoder) int64() int64 { return int64(d.uint64()) }

func (d *decoder) readPrefix() uint64 {
	var l uint64
	for shift := uint64(0); ; shift += 7 {
		b := d.buf[d.offset]
		d.offset++
		l |= uint64(b&0x7F) << shift
		if b < 0x80 {
			break
		}
	}
	return l
}

func (d *decoder) readBuf(l uint) []byte {
	return nil
}

func (d *decoder) readString(v reflect.Value) {
	l := int(d.readPrefix())
	if l > 0 {
		v.SetString(string(d.buf[d.offset : d.offset+l]))
		d.offset += l
	}
}

func (d *decoder) readArray(v reflect.Value) {
	l := v.Len()

	if k := v.Type().Elem().Kind(); k == reflect.Uint8 {
		v.SetBytes(d.buf[d.offset : d.offset+l])
		d.offset += l
	}

	for i := 0; i < l; i++ {
		d.read(v.Index(i))
	}
}

func (d *decoder) readSlice(v reflect.Value) {
	l := int(d.readPrefix())

	if k := v.Type().Elem().Kind(); k == reflect.Uint8 {
		v.SetBytes(d.buf[d.offset : d.offset+l])
		d.offset += l
	}

	tv := reflect.MakeSlice(v.Type(), l, l)
	for i := 0; i < l; i++ {
		d.read(tv.Index(i))
	}
	v.Set(tv)
}
