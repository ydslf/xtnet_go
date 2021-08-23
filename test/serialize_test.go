package test

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"testing"
	"unsafe"
)

type TestStructTobytes1 struct {
	Data   int64
	DArray [3]int32
	Data1  int32
}

func TestFloat(t *testing.T) {
	var f1 float32 = 1.0
	fmt.Printf("%#x\n", f1)

	var i1 int32 = 1
	fmt.Printf("%#x\n", i1)

	i2 := int32(f1)
	fmt.Printf("%v\n", i2)

	i3 := *(*int32)(unsafe.Pointer(&f1))
	fmt.Printf("%v\n", i3)
}

func TestBuffer(t *testing.T) {
	buf := &bytes.Buffer{}
	buf.WriteString("abc")
}

func TestStringSize(t *testing.T) {
	str := "a"
	for i := 0; i < 128; i++ {
		str += "a"
	}

	var s uint64 = 0
	l := uint64(len(str))
	t1 := l
	for t1 >= 0x80 {
		t1 >>= 7
		s++
	}
	s++
	s += l
	fmt.Println(s)
}

func TestGoBinary(t *testing.T) {
	valueEn := TestStructTobytes1{
		Data:  123,
		Data1: 789,
	}
	valueEn.DArray[0] = 4
	valueEn.DArray[1] = 5
	valueEn.DArray[2] = 6

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, valueEn)

	valueDe := &TestStructTobytes1{}
	buf1 := bytes.NewReader(buf.Bytes())
	binary.Read(buf1, binary.BigEndian, valueDe)
	fmt.Printf("%v\n", valueDe)
}

func TestGob(t *testing.T) {
	valueEn := TestStructTobytes1{
		Data:  123,
		Data1: 789,
	}
	valueEn.DArray[0] = 4
	valueEn.DArray[1] = 5
	valueEn.DArray[2] = 6

	buf := new(bytes.Buffer)
	encoder := gob.NewEncoder(buf)
	encoder.Encode(valueEn)

	decoder := gob.NewDecoder(bytes.NewReader(buf.Bytes()))
	valueDe := &TestStructTobytes1{}
	decoder.Decode(valueDe)
	fmt.Printf("%v\n", valueDe)
}

func TestJson(t *testing.T) {
	valueEn := TestStructTobytes1{
		Data:  123,
		Data1: 789,
	}
	valueEn.DArray[0] = 4
	valueEn.DArray[1] = 5
	valueEn.DArray[2] = 6

	b, _ := json.Marshal(valueEn)
	fmt.Println(b)
	fmt.Println(string(b))

	valueDe := &TestStructTobytes1{}
	json.Unmarshal(b, valueDe)
	fmt.Printf("%v\n", valueDe)
}

func BenchmarkGoBinary(b *testing.B) {
	valueEn := TestStructTobytes1{
		Data:  123,
		Data1: 789,
	}
	valueEn.DArray[0] = 4
	valueEn.DArray[1] = 5
	valueEn.DArray[2] = 6

	for i := 0; i < b.N; i++ {
		buf := new(bytes.Buffer)
		binary.Write(buf, binary.BigEndian, valueEn)

		valueDe := &TestStructTobytes1{}
		buf1 := bytes.NewReader(buf.Bytes())
		binary.Read(buf1, binary.BigEndian, valueDe)
	}
}

func BenchmarkGob(b *testing.B) {
	valueEn := TestStructTobytes1{
		Data:  123,
		Data1: 789,
	}
	valueEn.DArray[0] = 4
	valueEn.DArray[1] = 5
	valueEn.DArray[2] = 6

	for i := 0; i < b.N; i++ {
		buf := new(bytes.Buffer)
		encoder := gob.NewEncoder(buf)
		encoder.Encode(valueEn)

		decoder := gob.NewDecoder(bytes.NewReader(buf.Bytes()))
		valueDe := &TestStructTobytes1{}
		decoder.Decode(valueDe)
	}
}

func BenchmarkJson(b *testing.B) {
	valueEn := TestStructTobytes1{
		Data:  123,
		Data1: 789,
	}
	valueEn.DArray[0] = 4
	valueEn.DArray[1] = 5
	valueEn.DArray[2] = 6

	for i := 0; i < b.N; i++ {
		b, _ := json.Marshal(valueEn)

		valueDe := &TestStructTobytes1{}
		json.Unmarshal(b, valueDe)
	}
}
