package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
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