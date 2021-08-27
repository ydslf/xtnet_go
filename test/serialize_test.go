package test

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"testing"
	"unsafe"
	xt_encoding "xtnet/encoding"
	"xtnet/util"
)

type TestStructTobytes1 struct {
	Data   int64
	DArray [3]int32
	Data1  int32
}

type Person struct {
	Name   string
	Age    uint8
	Height float64
}

type Group struct {
	Name    string
	Members []Person
}

type A struct {
	Name     string
	BirthDay int64
	Phone    string
	Siblings int64
	Spouse   bool
	Money    float64
}
type B struct {
	BirthDay int64
	Siblings int64
	Spouse   bool
	Money    float64
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

func TestXtEncode(t *testing.T) {
	p := Group{
		Name: "test",
		Members: []Person{
			{
				Name:   "join",
				Age:    21,
				Height: 5.9,
			},
			{
				Name:   "Tom",
				Age:    23,
				Height: 5.8,
			},
			{
				Name:   "Alan",
				Age:    24,
				Height: 6,
			},
		},
	}

	fmt.Println(p)
	buf, _ := xt_encoding.Encode(p)
	fmt.Println(buf)

	pd := &Group{}
	xt_encoding.Decode(buf, pd)
	fmt.Println(pd)

	fmt.Println("-------------------------------------------")

	fmt.Println(p)
	buf1 := util.NewBuffer()
	xt_encoding.Encode2Buf(buf1, p)
	fmt.Println(buf1.GetReadData())

	pd1 := &Group{}
	xt_encoding.Decode(buf1.GetReadData(), pd1)
	fmt.Println(pd)
}

func TestGoBinaryArray(t *testing.T) {
	var enInt [3]int32
	enInt[0] = 1
	enInt[1] = 2
	enInt[2] = 3
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, enInt)

	var DeInt [3]int32
	reader := bytes.NewReader(buf.Bytes())
	binary.Read(reader, binary.BigEndian, &DeInt)
	fmt.Printf("%v\n", DeInt)
}

func TestGoBinarySlice(t *testing.T) {
	var enInt []int32
	enInt = append(enInt, 1, 2, 3)
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, enInt)

	var DeInt []int32
	reader := bytes.NewReader(buf.Bytes())
	binary.Read(reader, binary.BigEndian, &DeInt)
	fmt.Printf("%v\n", DeInt)
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

func BenchmarkEncode(b *testing.B) {
	p := Group{
		Name: "test",
		Members: []Person{
			{
				Name:   "John",
				Age:    21,
				Height: 5.9,
			},
			{
				Name:   "Tom",
				Age:    23,
				Height: 5.8,
			},
			{
				Name:   "Alan",
				Age:    24,
				Height: 6,
			},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		xt_encoding.Encode(p)
	}
}

func BenchmarkEncode2Buf(b *testing.B) {
	p := Group{
		Name: "test",
		Members: []Person{
			{
				Name:   "John",
				Age:    21,
				Height: 5.9,
			},
			{
				Name:   "Tom",
				Age:    23,
				Height: 5.8,
			},
			{
				Name:   "Alan",
				Age:    24,
				Height: 6,
			},
		},
	}

	buf := util.NewBuffer()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		xt_encoding.Encode2Buf(buf, p)
	}
}

func BenchmarkDecode(b *testing.B) {
	p := Group{
		Name: "test",
		Members: []Person{
			{
				Name:   "join",
				Age:    21,
				Height: 5.9,
			},
			{
				Name:   "Tom",
				Age:    23,
				Height: 5.8,
			},
			{
				Name:   "Alan",
				Age:    24,
				Height: 6,
			},
		},
	}

	buf, _ := xt_encoding.Encode(p)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		pd := &Group{}
		xt_encoding.Decode(buf, pd)
	}
}

func BenchmarkEncode1(b *testing.B) {
	p := A{
		Name:     "abcdefghijklmno",
		BirthDay: 123,
		Phone:    "1234567890",
		Siblings: 456,
		Spouse:   true,
		Money:    123.456,
	}

	for i := 0; i < b.N; i++ {
		xt_encoding.Encode(p)
	}
}

func BenchmarkEncodeFixStruct(b *testing.B) {
	p := B{
		BirthDay: 123,
		Siblings: 456,
		Spouse:   true,
		Money:    123.456,
	}

	for i := 0; i < b.N; i++ {
		xt_encoding.Encode(p)
	}
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
