package main

import (
	"bufio"
	//"fmt"
	"os"
	"strconv"
	"strings"
	//"testing"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

//func main2() {
// result := testing.Benchmark(run)
//fmt.Printf("Time taken: %s \n", result)
//fmt.Printf("Number of bytes allocated: %d \n", result.Bytes)
//fmt.Printf("Memory allocations: %d \n", result.MemAllocs)
//}

// func run(b *testing.B) {
func main() {
	input, err := os.Open("../measurements.txt")
	check(err)
	defer input.Close()

	counts := make(map[string]int)
	avgs := make(map[string]float32)

	scanner := bufio.NewScanner(input)
	//total := 0
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), ";")
		city := line[0]
		temp64, _ := strconv.ParseFloat(line[1], 32)
		temp := float32(temp64)

		count := counts[city] + 1
		tempAvg := float32(avgs[city] * float32(count-1))

		avgs[city] = float32(tempAvg+temp) / float32(count)
		counts[city] = count

		//total += 1
		//if total%25000000 == 0 {
		//	fmt.Println(total)
		//}
	}

	//for city, avg := range avgs {
	//	fmt.Println(city, avg)
	//}
}
