package main

import (
	"testing"
)

func BenchmarkTest1BRC(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Run1BRC(true)
	}
}

func BenchmarkReal1BRC(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Run1BRC(false)
	}
}
