package main

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

var m sync.Mutex
var set = make(map[int]bool)

func printOnce(num int) {
	m.Lock()
	defer m.Unlock()
	if _, exist := set[num]; !exist {
		fmt.Println(num)
	}
	set[num] = true

}

func TestPrintOnce(t *testing.T) {
	for i := 0; i < 10; i++ {
		go printOnce(100)
	}
	time.Sleep(2 * time.Second)
}
