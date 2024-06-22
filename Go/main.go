package main

import (
	"bufio"
	"fmt"
	"sync"

	// "math"
	"os"
	"runtime/pprof"
	"sort"
	"strconv"
	"strings"
)

type city struct {
	count int
	total float64
	min   float64
	max   float64
}

func (c city) process(in float64) {
	c.count += 1
	c.total += in
	if in < c.min {
		c.min = in
		return
	} else if in > c.max {
		c.max = in
	}
}

func (c city) getAvg() float64 {
	return c.total / float64(c.count)
}

type mapHandler struct {
	mapping map[string]city
}

func (handler mapHandler) process(name string, in string) {
	c, exist := handler.mapping[name]
	if !exist {
		c = city{}
	}

	float, err := strconv.ParseFloat(in, 64)
	check(err)

	c.process(float)
}

func (handler mapHandler) getSortedKeys() []string {
	mapping := handler.mapping
	keys := make([]string, 0, len(mapping))
	for k := range mapping {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func (handler mapHandler) getCity(name string) city {
	return handler.mapping[name]
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	file, err := os.Create("1BRC.prof")
	check(err)
	pprof.StartCPUProfile(file)
	Run1BRC()
	pprof.StopCPUProfile()
}

const ChunkSize = 2048

func Run1BRC() {
	input, err := os.Open("../measurements.txt")
	check(err)
	defer input.Close()

	handler := mapHandler{}

	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), ";")
		city, temp := line[0], line[1]
		handler.process(city, temp)
	}

	keys := handler.getSortedKeys()

	for k := range keys {
		c := handler.getCity(keys[k])

		fmt.Println(fmt.Sprintf("%s %.1f %.1f %.1f", keys[k], c.min, c.getAvg(), c.max))
	}
}

// For when I implement goroutines. Presently just want to see if using structs helps anything
func ProcessRows(row string, wg *sync.WaitGroup) {
	defer wg.Done()
}
