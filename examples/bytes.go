package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"
	"unsafe"
)

type TestStructTobytes struct {
	Data   int64
	DSlice []int32
	Data1  int32
}

type TestStructTobytes1 struct {
	Data   int64
	DArray [3]int32
	Data1  int32
}

type SliceMock struct {
	addr uintptr
	len  int
	cap  int
}

func testBytes1() {
	var testStruct = &TestStructTobytes{
		Data: 1234,
	}
	Len := unsafe.Sizeof(*testStruct)
	testBytes := &SliceMock{
		addr: uintptr(unsafe.Pointer(testStruct)),
		cap:  int(Len),
		len:  int(Len),
	}

	fmt.Println(testStruct)
	fmt.Printf("%p\n", testStruct)
	fmt.Println(unsafe.Pointer(testStruct))

	data := *(*[]byte)(unsafe.Pointer(testBytes))
	fmt.Println("[]byte is : ", data)
	var ptestStruct = *(**TestStructTobytes)(unsafe.Pointer(&data))
	fmt.Println("ptestStruct.data is : ", ptestStruct.Data)

	test1 := []byte{111, 0, 1, 2}
	test2 := string(test1)
	test3, _ := strconv.Atoi(test2)
	fmt.Println("Atoi: ", test3)

	var a int32
	fmt.Println(a)
	fmt.Println(unsafe.Sizeof(a))
	fmt.Printf("%T\n", unsafe.Sizeof(a))

	var s1 []byte
	fmt.Println(s1)
	fmt.Println(nil)

	var s2 []byte = nil
	fmt.Println(s2)
}

func testBytesSlice() {
	buf := new(bytes.Buffer)
	valueSource := []int32{1, 2, 3, 4, 5, 6, 7, 8}
	valueEn := valueSource[:5]
	e1 := binary.Write(buf, binary.BigEndian, valueEn)
	fmt.Println(e1)
	valueEnDe := make([]int32, 3)
	buf1 := bytes.NewReader(buf.Bytes())
	e2 := binary.Read(buf1, binary.BigEndian, valueEnDe)
	fmt.Printf("%v\n", valueEnDe)
	fmt.Println(e2)
}

func testBytesStruct() {
	buf := new(bytes.Buffer)
	valueEn := &TestStructTobytes{
		Data:  123,
		Data1: 789,
	}
	valueEn.DSlice = append(valueEn.DSlice, 4)
	valueEn.DSlice = append(valueEn.DSlice, 5)
	valueEn.DSlice = append(valueEn.DSlice, 6)

	e1 := binary.Write(buf, binary.BigEndian, valueEn)
	fmt.Println(e1)
	valueEnDe := &TestStructTobytes{}
	buf1 := bytes.NewReader(buf.Bytes())
	e2 := binary.Read(buf1, binary.BigEndian, valueEnDe)
	fmt.Printf("%v\n", valueEnDe)
	fmt.Println(e2)
}

func testBytesStruct1() {
	buf := new(bytes.Buffer)
	valueEn := TestStructTobytes1{
		Data:  123,
		Data1: 789,
	}
	valueEn.DArray[0] = 4
	valueEn.DArray[1] = 5
	valueEn.DArray[2] = 6

	e1 := binary.Write(buf, binary.BigEndian, valueEn)
	fmt.Println(e1)
	valueEnDe := &TestStructTobytes1{}
	buf1 := bytes.NewReader(buf.Bytes())
	e2 := binary.Read(buf1, binary.BigEndian, valueEnDe)
	fmt.Printf("%v\n", valueEnDe)
	fmt.Println(e2)
}

func main() {
	//testBytes1()
	//testBytesSlice()
	//testBytesStruct()
	testBytesStruct1()
}
