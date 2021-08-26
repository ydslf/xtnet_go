package test

import (
	"fmt"
	"reflect"
	"testing"
)

func TestReflect(t *testing.T) {
	var a = 56
	value := reflect.ValueOf(a)
	fmt.Println(value.Interface())
	fmt.Println(value.Interface().(int))

	str := "你好"
	fmt.Println(len(str))

	var b uint8 = 0xFF
	c := int8(b)
	var d = int(int8(b))
	fmt.Println(b)
	fmt.Println(c)
	fmt.Println(d)
}

type TestStruct struct {
	Data  int64
	Data1 int32
}

func myTestFun(obj interface{}) {
	fmt.Printf("%p\n", &obj)
	//v := reflect.ValueOf(obj)
	//v.Field(0).SetInt(1)
	//obj.Data = 1
	fmt.Println(obj)
}
func TestInterface(t *testing.T) {
	obj := TestStruct{
		Data:  123,
		Data1: 789,
	}

	fmt.Println(obj)
	fmt.Printf("%p\n", &obj)
	myTestFun(obj)
	fmt.Println(obj)
}

func tt(a interface{}) (int, string) {
	return a.(int), "abc"
}
func TestTypeChange(t *testing.T) {

	var b1 int8 = 1
	fmt.Println(b1)
	fmt.Println(b1 & 0)
	fmt.Println(b1 | 0)
	fmt.Println("------------------")

	var a interface{} = 10

	fmt.Println(tt(a))

	v, ok := a.(int)
	fmt.Println(v, ok)

	v2, ok := a.(float32)
	fmt.Println(v2, ok)
}
