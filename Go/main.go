package main

import (
	"bufio"
	"fmt"
	// "math"
	"os"
	"sort"
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
	totals := make(map[string]float64)
	mins := make(map[string]float64)
	maxs := make(map[string]float64)

	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), ";")
		city := line[0]
		temp, _ := strconv.ParseFloat(line[1], 32)

		totals[city] += temp
		counts[city] += 1

		if temp < mins[city] {
			mins[city] = temp
			continue
		} else if temp > maxs[city] {
			maxs[city] = temp
		}

	}

	keys := make([]string, 0, len(counts))
	for k := range counts {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		avg := totals[k] / float64(counts[k])
		fmt.Println(fmt.Sprintf("%s %.1f %.1f %.1f", k, mins[k], avg, maxs[k]))
		// fmt.Println(k, mins[k], math.Floor(avg*10)/10, maxs[k])
	}
}
