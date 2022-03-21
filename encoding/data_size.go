package encoding

import (
	"errors"
	"reflect"
)

var (
	ECErrUnSupportType = errors.New("type is not support")
)

func dataSize(v reflect.Value) (int, int8, error) {
	var n = 0
	var fixSize int8 = 1

	switch v.Kind() {
	case reflect.Bool,
		reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:
		n += int(v.Type().Size())

	case reflect.String:
		l := v.Len()
		n += calcPrefixSize(uint64(l))
		n += l
		fixSize = 0

	case reflect.Array:
		size, fix, err := calcArraySize(v)
		if err != nil {
			return 0, 0, err
		}
		n += size
		fixSize &= fix

	case reflect.Slice:
		size, _, err := calcSliceSize(v)
		if err != nil {
			return 0, 0, err
		}
		n += size
		fixSize = 0

	case reflect.Struct:
		size, fix, err := calcStructSize(v)
		if err != nil {
			return 0, 0, err
		}
		n += size
		fixSize &= fix

	default:
		return 0, 0, ECErrUnSupportType
	}

	return n, fixSize, nil
}

func calcPrefixSize(l uint64) (n int) {
	for {
		n++
		l >>= 7
		if l == 0 {
			break
		}
	}
	return n
}

func calcArraySize(v reflect.Value) (int, int8, error) {
	var n = 0
	var fixSize int8 = 1

	l := v.Len()
	elem := v.Type().Elem()
	switch elem.Kind() {
	case reflect.Bool,
		reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:
		n += l * int(elem.Size())

	default:
		for i := 0; i < l; i++ {
			size, fix, err := dataSize(v.Index(i))
			if err != nil {
				return 0, 0, err
			}
			n += size
			fixSize &= fix
		}
	}
	return n, fixSize, nil
}

func calcSliceSize(v reflect.Value) (int, int8, error) {
	var n = 0
	l := v.Len()

	n += calcPrefixSize(uint64(l))

	elem := v.Type().Elem()
	switch elem.Kind() {
	case reflect.Bool,
		reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:
		n += l * int(elem.Size())

	default:
		for i := 0; i < l; i++ {
			size, _, err := dataSize(v.Index(i))
			if err != nil {
				return 0, 0, err
			}
			n += size
		}
	}

	return n, 0, nil
}
