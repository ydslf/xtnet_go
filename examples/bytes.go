package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
	"strconv"
	"unsafe"
)

type TestStructTobytes11 struct {
	Data int64
}

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

type TestStructTobytesPrivate struct {
	data  int64
	Data1 int32
}

type TestStructTobytesReflect struct {
	data  int64 `key1:value1 key2:value2`
	Data1 int32
	Data2 string
	Data3 [5]int32
	_     int32
	_     int32
}

type TestStructTobytesSkip struct {
	Data  int64
	Data1 int32
}

type SliceMock struct {
	addr uintptr
	len  int
	cap  int
}

type StringMock struct {
	str uintptr
	len int
}

func testBytes1() {
	var testStruct = &TestStructTobytes11{
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
	var ptestStruct = *(**TestStructTobytes11)(unsafe.Pointer(&data))
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
	fmt.Println("==========================================================")
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
	fmt.Println("==========================================================")
}

func testBytesStructPrivate() {
	buf := new(bytes.Buffer)
	valueEn := TestStructTobytesPrivate{
		data:  123,
		Data1: 789,
	}

	e1 := binary.Write(buf, binary.BigEndian, valueEn)
	fmt.Println(e1)
	valueEnDe := &TestStructTobytesPrivate{}
	buf1 := bytes.NewReader(buf.Bytes())
	e2 := binary.Read(buf1, binary.BigEndian, valueEnDe)
	fmt.Printf("%v\n", valueEnDe)
	fmt.Println(e2)
}

func testBytesStructSkip() {
	buf := new(bytes.Buffer)
	valueEn := TestStructTobytesSkip{
		Data:  123,
		Data1: 789,
	}

	e1 := binary.Write(buf, binary.BigEndian, valueEn)
	fmt.Println(e1)
	valueEnDe := &TestStructTobytesSkip{}
	buf1 := bytes.NewReader(buf.Bytes())
	e2 := binary.Read(buf1, binary.BigEndian, valueEnDe)
	fmt.Printf("%v\n", valueEnDe)
	fmt.Println(e2)
}

func testReflect() {
	fmt.Println("==========================================================")
	a := &TestStructTobytesReflect{
		data:  1,
		Data1: 2,
	}
	fmt.Println(a)
	v := reflect.ValueOf(a)
	t := v.Type()
	fmt.Println(v)
	fmt.Println(t)
	fmt.Println(t.String())

	v = v.Elem()
	t = v.Type()
	fmt.Println(v)
	fmt.Println(t)
	for i := 0; i < t.NumField(); i++ {
		vv := v.Field(i)
		sf := t.Field(i)
		fmt.Println(vv)
		fmt.Println(sf)
		fmt.Println(sf.Type)
		fmt.Println(sf.Tag)
		fmt.Println(vv.CanSet())
		fmt.Println("------------------------")
	}

	var teststr string = "abc"
	v = reflect.ValueOf(&teststr)
	t = v.Type()
	v1 := v.Elem()
	v2 := reflect.ValueOf("abc")
	fmt.Printf("%p\n", &teststr)
	fmt.Println(v)
	fmt.Println(t)
	fmt.Println(v.Kind())
	fmt.Println(v.CanSet())
	v.Elem().SetString("def")
	fmt.Println(v1)
	fmt.Println(v2)
	fmt.Println(teststr)
	fmt.Printf("%p\n", &teststr)

	var testint int32 = 123
	v = reflect.ValueOf(&testint)
	fmt.Println(v.Elem().CanSet())

	fmt.Println("==========================================================")
}

func testStr() {
	fmt.Println("==========================================================")
	var teststr string = "abc"
	teststr1 := teststr
	fmt.Println(teststr)
	fmt.Println(teststr1)
	teststr = "def"
	fmt.Println(teststr)
	fmt.Println(teststr1)

	fmt.Println("------------------------")
	var teststr2 string = "abc"
	teststr3 := teststr2
	fmt.Println(teststr2)
	fmt.Println(teststr3)
	v := reflect.ValueOf(&teststr2)
	v.Elem().SetString("aaa")
	fmt.Println(teststr2)
	fmt.Println(teststr3)

	fmt.Println("------------------------")
	var data [10]byte
	data[0] = 'T'
	data[1] = 'E'
	var str string = string(data[:])
	fmt.Println(str)

	fmt.Println("------------------------")
	slice1 := []byte("abc")
	var str1 string = "def"
	pbytes := (*reflect.SliceHeader)(unsafe.Pointer(&slice1))
	pstring := (*reflect.StringHeader)(unsafe.Pointer(&str1))

	pstring.Data = pbytes.Data
	pstring.Len = pbytes.Len

	str2 := str1
	fmt.Println(str1)
	fmt.Println(str2)

	slice1[0] = 'x'
	fmt.Println(str1)
	fmt.Println(str2)

	str2 = "123"

	fmt.Println(str1)
	fmt.Println(str2)

	fmt.Println("==========================================================")
}

func main() {
	//testBytes1()
	//testBytesSlice()
	//testBytesStruct()
	//testBytesStruct1()
	//testBytesStructPrivate()
	//testBytesStructSkip()
	//testReflect()
	testStr()
}
