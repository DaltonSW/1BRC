package main

// NOTE: Looking back on this, I realize it actually didn't do the min/max checks
// ... but it was so slow that I'm not sure it's worth going back to update

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

const KBs = 1024
const MBs = 1024 * KBs

func main() {
	start := time.Now()
	Run1BRC()
	elapsed := time.Since(start)
	fmt.Printf("Processing took %s\n", elapsed)
}

func Run1BRC() {
	input, err := os.Open("../measurements.txt")
	check(err)
	defer input.Close()

	counts := make(map[string]int)
	avgs := make(map[string]float32)

	scanner := bufio.NewScanner(input)

	for scanner.Scan() {
		line := strings.Split(scanner.Text(), ";")
		city := line[0]
		temp64, _ := strconv.ParseFloat(line[1], 32)
		temp := float32(temp64)

		count := counts[city] + 1
		tempAvg := float32(avgs[city] * float32(count-1))

		avgs[city] = float32(tempAvg+temp) / float32(count)
		counts[city] = count

	}
}
