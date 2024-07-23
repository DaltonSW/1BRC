package main

import (
	"math/rand"
	"testing"
)

const NumInts = 10
const RangeMin = -100
const RangeMax = 100

func MinMaxSetup() []int {
	nums := make([]int, NumInts)
	for i := range NumInts {
		nums[i] = rand.Intn(RangeMax-RangeMin+1) + RangeMin
	}
	return nums
}

func BenchmarkMinMaxWithFunc(b *testing.B) {
	nums := MinMaxSetup()

	for i := 0; i < b.N; i++ {
		minNum := -9999
		maxNum := 9999

		for _, n := range nums {
			minNum = min(minNum, n)
			maxNum = max(maxNum, n)
		}
	}

}

func BenchmarkMinMaxManual(b *testing.B) {
	nums := MinMaxSetup()

	for i := 0; i < b.N; i++ {
		minNum := -9999
		maxNum := 9999

		for _, n := range nums {
			if n < minNum {
				minNum = n
				break
			} else if n > maxNum {
				maxNum = n
			}
		}
	}
}

// func TestRunProgram(t *testing.T) {
// 	Run1BRC(true, 1*MBs)
// }

// func BenchmarkReal1BRC(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		Run1BRC(false, 2*MBs)
// 	}
// }
