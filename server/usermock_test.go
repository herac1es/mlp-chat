package server

import (
	"fmt"
	"testing"
)

func Benchmark_generate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		generate()
	}
}

func Test_generate(t *testing.T) {
	for i := 0; i < 1000; i++ {
		fmt.Println(generate())
	}
}
