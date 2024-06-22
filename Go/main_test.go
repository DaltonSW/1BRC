package main

import (
	"testing"
)

func Benchmark1BRC(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Run1BRC()
	}
}
