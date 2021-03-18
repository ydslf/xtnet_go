package main

import (
	"fmt"
	"strconv"
	"unsafe"
)

type TestStructTobytes struct {
	data int64
}
type SliceMock struct {
	addr uintptr
	len  int
	cap  int
}

func main() {
	var testStruct = &TestStructTobytes{1234}
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
	fmt.Println("ptestStruct.data is : ", ptestStruct.data)

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
