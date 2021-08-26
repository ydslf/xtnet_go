package encoding

import (
	"reflect"
	"sync"
)

type structCache struct {
	indexes []int
}

var structSize sync.Map
var structCaches sync.Map

func calcStructSize(v reflect.Value) (int, int8, error) {
	var n = 0
	t := v.Type()

	if size, find := structSize.Load(t); find {
		return size.(int), 1, nil
	}

	var fixSize int8 = 1
	cache, find := structCaches.Load(t)
	if find {
		sc := cache.(*structCache)
		for _, structIndex := range sc.indexes {
			size, fix, err := dataSize(v.Field(structIndex))
			if err != nil {
				return 0, 0, err
			}
			n += size
			fixSize &= fix
		}
	} else {
		sc := &structCache{}
		for i := 0; i < v.NumField(); i++ {
			if sf := t.Field(i); sf.Tag.Get("xt") != "-" {
				size, fix, err := dataSize(v.Field(i))
				if err != nil {
					return 0, 0, err
				}
				n += size
				fixSize &= fix
				sc.indexes = append(sc.indexes, i)
			}
		}
		structCaches.Store(t, sc)
	}

	if fixSize == 1 {
		structSize.Store(t, n)
	}

	return n, fixSize, nil
}

func (e *encoder) writeStruct(v reflect.Value) {
	t := v.Type()

	cache, find := structCaches.Load(t)
	if find {
		sc := cache.(*structCache)
		for _, structIndex := range sc.indexes {
			e.write(v.Field(structIndex))
		}
	} else {
		sc := &structCache{}
		for i := 0; i < v.NumField(); i++ {
			if sf := t.Field(i); sf.Tag.Get("xt") != "-" {
				e.write(v.Field(i))
				sc.indexes = append(sc.indexes, i)
			}
		}
		structCaches.Store(t, sc)
	}
}

func (d *decoder) readStruct(v reflect.Value) {
	t := v.Type()

	cache, find := structCaches.Load(t)
	if find {
		sc := cache.(*structCache)
		for _, structIndex := range sc.indexes {
			d.read(v.Field(structIndex))
		}
	} else {
		sc := &structCache{}
		for i := 0; i < v.NumField(); i++ {
			if sf := t.Field(i); sf.Tag.Get("xt") != "-" {
				d.read(v.Field(i))
				sc.indexes = append(sc.indexes, i)
			}
		}
		structCaches.Store(t, sc)
	}
}
