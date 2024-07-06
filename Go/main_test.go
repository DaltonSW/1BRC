package main

import (
	"testing"
)

func BenchmarkTest512MB(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Run1BRC(true, 512*MBs)
	}
}

func BenchmarkTest256MB(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Run1BRC(true, 256*MBs)
	}
}

func BenchmarkTest128MB(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Run1BRC(true, 128*MBs)
	}
}

func BenchmarkTest64MB(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Run1BRC(true, 64*MBs)
	}
}

func BenchmarkTest32MB(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Run1BRC(true, 32*MBs)
	}
}

func BenchmarkTest16MB(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Run1BRC(true, 16*MBs)
	}
}

func BenchmarkTest8MB(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Run1BRC(true, 8*MBs)
	}
}
func BenchmarkTest4MB(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Run1BRC(true, 4*MBs)
	}
}
func BenchmarkTest2MB(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Run1BRC(true, 2*MBs)
	}
}
func BenchmarkReal1BRC(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Run1BRC(false, 32*MBs)
	}
}
