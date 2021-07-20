package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"strings"
	"testing"
)

func TestCatString(t *testing.T) {
	fmt.Println("test")
}

func benchmark(b *testing.B, f func(string, int) string) {
	str := generateRandomString(10)
	for i := 0; i < b.N; i++ {
		f(str, 10000)
	}
}

func BenchmarkPlusCat(b *testing.B) {
	benchmark(b, plusCat)
}
func BenchmarkSprintfCat(b *testing.B) {
	benchmark(b, sprintfCat)
}
func BenchmarkBuilderCat(b *testing.B) {
	benchmark(b, builderCat)
}
func BenchmarkBufferCat(b *testing.B) {
	benchmark(b, bufferCat)
}
func BenchmarkByteCat(b *testing.B) {
	benchmark(b, byteCat)
}

func generateRandomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = byte(rand.Intn(128))
	}
	return string(b)
}

func plusCat(s string, times int) string {
	var ret string
	for i := 0; i < times; i++ {
		ret += s
	}
	return ret
}

func sprintfCat(s string, times int) string {
	var ret string
	for i := 0; i < times; i++ {
		ret = fmt.Sprintf("%s%s", ret, s)
	}
	return ret
}

func builderCat(s string, times int) string {
	var builder strings.Builder
	for i := 0; i < times; i++ {
		builder.WriteString(s)
	}
	return builder.String()
}

func bufferCat(s string, times int) string {
	buf := new(bytes.Buffer)
	for i := 0; i < times; i++ {
		buf.WriteString(s)
	}
	return buf.String()
}

func byteCat(s string, times int) string {
	buf := make([]byte, 0)
	for i := 0; i < times; i++ {
		buf = append(buf, s...)
	}
	return string(buf)
}
