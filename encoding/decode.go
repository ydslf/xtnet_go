package encoding

import (
	"encoding/binary"
	"errors"
	"math"
	"reflect"
)

var (
	ECErrMustPtrType = errors.New("decode interface must be ptr")
	ECErrOutOfRange  = errors.New("decode out of range")
)

func Decode(data []byte, obj interface{}) error {
	v := reflect.ValueOf(obj)

	if v.Kind() != reflect.Ptr {
		return ECErrMustPtrType
	}

	v = v.Elem()
	d := decoder{order: binary.BigEndian, buf: data}
	d.dataSize = len(data)
	err := d.read(v)

	return err
}

type decoder struct {
	order    binary.ByteOrder
	buf      []byte
	dataSize int
	offset   int
}

func (d *decoder) read(v reflect.Value) error {
	var err error

	switch v.Kind() {
	case reflect.Bool:
		err = d.readBool(v)

	case reflect.Int8:
		err = d.readInt8(v)
	case reflect.Int16:
		err = d.readInt16(v)
	case reflect.Int32:
		err = d.readInt32(v)
	case reflect.Int64:
		err = d.readInt64(v)

	case reflect.Uint8:
		err = d.readUint8(v)
	case reflect.Uint16:
		err = d.readUint16(v)
	case reflect.Uint32:
		err = d.readUint32(v)
	case reflect.Uint64:
		err = d.readUint64(v)

	case reflect.Float32:
		err = d.readFloat32(v)
	case reflect.Float64:
		err = d.readFloat64(v)

	case reflect.Complex64:
		err = d.readComplex64(v)
	case reflect.Complex128:
		err = d.readComplex128(v)

	case reflect.String:
		err = d.readString(v)

	case reflect.Array:
		err = d.readArray(v)

	case reflect.Slice:
		err = d.readSlice(v)

	case reflect.Struct:
		err = d.readStruct(v)
	}
	return err
}

func (d *decoder) checkSize(size int) bool {
	return size <= (d.dataSize - d.offset)
}

func (d *decoder) readBool(v reflect.Value) error {
	if !d.checkSize(1) {
		return ECErrOutOfRange
	}
	x := d.buf[d.offset]
	d.offset++
	v.SetBool(x != 0)
	return nil
}

func (d *decoder) readUint8(v reflect.Value) error {
	if !d.checkSize(1) {
		return ECErrOutOfRange
	}
	x := d.buf[d.offset]
	d.offset++
	v.SetUint(uint64(x))
	return nil
}

func (d *decoder) readUint16(v reflect.Value) error {
	if !d.checkSize(2) {
		return ECErrOutOfRange
	}
	x := d.order.Uint16(d.buf[d.offset : d.offset+2])
	d.offset += 2
	v.SetUint(uint64(x))
	return nil
}

func (d *decoder) readUint32(v reflect.Value) error {
	if !d.checkSize(4) {
		return ECErrOutOfRange
	}
	x := d.order.Uint32(d.buf[d.offset : d.offset+4])
	d.offset += 4
	v.SetUint(uint64(x))
	return nil
}

func (d *decoder) readUint64(v reflect.Value) error {
	if !d.checkSize(8) {
		return ECErrOutOfRange
	}
	x := d.order.Uint64(d.buf[d.offset : d.offset+8])
	d.offset += 8
	v.SetUint(x)
	return nil
}

func (d *decoder) readInt8(v reflect.Value) error {
	if !d.checkSize(1) {
		return ECErrOutOfRange
	}
	x := d.buf[d.offset]
	d.offset++
	v.SetInt(int64(x))
	return nil
}

func (d *decoder) readInt16(v reflect.Value) error {
	if !d.checkSize(2) {
		return ECErrOutOfRange
	}
	x := d.order.Uint16(d.buf[d.offset : d.offset+2])
	d.offset += 2
	v.SetInt(int64(x))
	return nil
}

func (d *decoder) readInt32(v reflect.Value) error {
	if !d.checkSize(4) {
		return ECErrOutOfRange
	}
	x := d.order.Uint32(d.buf[d.offset : d.offset+4])
	d.offset += 4
	v.SetInt(int64(x))
	return nil
}

func (d *decoder) readInt64(v reflect.Value) error {
	if !d.checkSize(8) {
		return ECErrOutOfRange
	}
	x := d.order.Uint16(d.buf[d.offset : d.offset+8])
	d.offset += 8
	v.SetInt(int64(x))
	return nil
}

func (d *decoder) readFloat32(v reflect.Value) error {
	if !d.checkSize(4) {
		return ECErrOutOfRange
	}
	x := d.order.Uint32(d.buf[d.offset : d.offset+4])
	d.offset += 4
	v.SetFloat(float64(math.Float32frombits(x)))
	return nil
}

func (d *decoder) readFloat64(v reflect.Value) error {
	if !d.checkSize(8) {
		return ECErrOutOfRange
	}
	x := d.order.Uint64(d.buf[d.offset : d.offset+8])
	d.offset += 8
	v.SetFloat(math.Float64frombits(x))
	return nil
}

func (d *decoder) readComplex64(v reflect.Value) error {
	if !d.checkSize(8) {
		return ECErrOutOfRange
	}
	x1 := d.order.Uint32(d.buf[d.offset : d.offset+4])
	x2 := d.order.Uint32(d.buf[d.offset+4 : d.offset+8])
	d.offset += 8

	v.SetComplex(complex(
		float64(math.Float32frombits(x1)),
		float64(math.Float32frombits(x2)),
	))

	return nil
}

func (d *decoder) readComplex128(v reflect.Value) error {
	if !d.checkSize(16) {
		return ECErrOutOfRange
	}
	x1 := d.order.Uint64(d.buf[d.offset : d.offset+8])
	x2 := d.order.Uint64(d.buf[d.offset+8 : d.offset+16])
	d.offset += 16

	v.SetComplex(complex(
		math.Float64frombits(x1),
		math.Float64frombits(x2),
	))

	return nil
}

func (d *decoder) readPrefix() (uint64, error) {
	var l uint64
	for shift := uint64(0); ; shift += 7 {
		if !d.checkSize(1) {
			return 0, ECErrOutOfRange
		}
		b := d.buf[d.offset]
		d.offset++
		l |= uint64(b&0x7F) << shift
		if b < 0x80 {
			break
		}
	}
	return l, nil
}

func (d *decoder) readString(v reflect.Value) error {
	l, err := d.readPrefix()
	if err != nil {
		return err
	}

	l1 := int(l)
	if l1 > 0 {
		if !d.checkSize(l1) {
			return ECErrOutOfRange
		}
		v.SetString(string(d.buf[d.offset : d.offset+l1]))
		d.offset += l1
	}
	return nil
}

func (d *decoder) readArray(v reflect.Value) error {
	l := v.Len()

	if k := v.Type().Elem().Kind(); k == reflect.Uint8 {
		if !d.checkSize(l) {
			return ECErrOutOfRange
		}
		v.SetBytes(d.buf[d.offset : d.offset+l])
		d.offset += l
		return nil
	}

	for i := 0; i < l; i++ {
		err := d.read(v.Index(i))
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *decoder) readSlice(v reflect.Value) error {
	l, err := d.readPrefix()
	if err != nil {
		return err
	}

	l1 := int(l)

	if k := v.Type().Elem().Kind(); k == reflect.Uint8 {
		if !d.checkSize(l1) {
			return ECErrOutOfRange
		}
		v.SetBytes(d.buf[d.offset : d.offset+l1])
		d.offset += l1
	}

	tv := reflect.MakeSlice(v.Type(), l1, l1)
	for i := 0; i < l1; i++ {
		err = d.read(tv.Index(i))
		if err != nil {
			return err
		}
	}
	v.Set(tv)

	return nil
}
