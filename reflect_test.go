package main

import (
	"fmt"
	"reflect"
	"testing"
)

type S struct {
	x int
	y string
}

func TestReflect(t *testing.T) {
	a := 0
	s := S{
		x: 1,
		y: "abc",
	}
	f := func(x int, y int) int {
		return x - y
	}
	testReflect(t, a)
	testReflect(t, s)
	testReflect(t, f)
}

func testReflect(t *testing.T, a interface{}) {
	fmt.Println(reflect.TypeOf(a))
	fmt.Println(reflect.ValueOf(a))
}
