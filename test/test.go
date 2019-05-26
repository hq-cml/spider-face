package main

/*
 * 一些高级用法的测试
 */
import (
	"fmt"
	"reflect"
)

func test1() {
	var i interface{}
	var a interface{}
	a = 10
	i = &a

	fmt.Println(reflect.ValueOf(a))
	fmt.Println(reflect.ValueOf(i))

	fmt.Println(reflect.Indirect(reflect.ValueOf(a)))
	fmt.Println(reflect.Indirect(reflect.ValueOf(i)))
}

func test2() {
	a := 10
	fmt.Println(reflect.TypeOf(a))

	b := reflect.New(reflect.TypeOf(a))
	fmt.Println(b)
	bElem := b.Elem()
	fmt.Println(bElem)
	fmt.Println(bElem.Type())
	fmt.Println(b.CanSet())
	fmt.Println(bElem.CanSet())

	bElem.SetInt(88)
	fmt.Println(bElem.Int())
	fmt.Println(bElem.Interface())
}

func main() {
	test2()
}