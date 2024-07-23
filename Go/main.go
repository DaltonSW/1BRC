package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"sync"
	"time"

	"os"
	"runtime/pprof"
	"sort"
	"strconv"
	"strings"
)

// ISSUE: Worst Profiling Offenders are marked with issues
//	- ProcessChunk -- strings.Split (18.4%)
//	- Run1BRC -- strings.Split (5.1%)

// IDEA: Figure out how to parse the chunk in such a way that the goroutines can handle splitting

// IDEA: Multiply string by 10 to convert float into int to make math better for CPU

// IDEA: See if converting int to string and manually slotting a period in is faster than int -> float and string format

// IDEA: Parse the string backwards

// IDEA: Don't use structs and try normal vars and funcs

// IDEA: Try int32 instead of int64

const KBs = 1024
const MBs = 1024 * KBs

func main() {
	file, err := os.Create("1BRC.prof")
	check(err)
	pprof.StartCPUProfile(file)
	defer pprof.StopCPUProfile()

	countPtr := flag.Int("count", 10, "Number of tests to run and average")
	flag.Parse()

	testCount := *countPtr

	start := time.Now()
	for i := 0; i < testCount; i++ {
		Run1BRC(false, 8*MBs)
	}
	elapsed := time.Since(start)
	average := elapsed.Seconds() / float64(testCount)
	fmt.Printf("Processing %d tests.\n", testCount)
	fmt.Printf("Took a total of %s\n", elapsed)
	fmt.Printf("Took an average of %.3fs\n", average)
}

func Run1BRC(test bool, bufferSize int) {
	var input *os.File
	var err error

	if test {
		input, err = os.Open("../test_measurements.txt")
	} else {
		input, err = os.Open("../measurements.txt")
	}
	check(err)
	defer input.Close()

	lineBuffer := make([]byte, bufferSize)
	fileReader := bufio.NewReader(input)

	var wg sync.WaitGroup
	handler := mapHandler{mapping: make(map[string]*city)}
	remainder := ""

	for {
		num, err := fileReader.Read(lineBuffer)
		if num == 0 {
			if err == io.EOF {
				break
			}
			check(err)
			break
		}

		chunk := remainder + string(lineBuffer[:num])
		lines := strings.Split(chunk, "\n")

		remainder = lines[len(lines)-1]
		lines = lines[:len(lines)-1]

		wg.Add(1)
		go func(lines []string) {
			defer wg.Done()
			ProcessChunk(&handler, lines)
		}(lines)
	}

	if remainder != "" {
		ProcessChunk(&handler, []string{remainder})
	}

	wg.Wait()

	keys := handler.getSortedKeys()

	for k := range keys {
		c := handler.getCity(keys[k])

		fmt.Println(fmt.Sprintf("%s=%.1f/%.1f/%.1f", keys[k], c.min, c.getAvg(), c.max))
	}
}

type city struct {
	count int
	total float64
	min   float64
	max   float64
	mu    sync.RWMutex
}

func NewCity() *city {
	c := city{}

	c.max = -1e99
	c.min = 1e99

	return &c
}

func (c *city) process(in float64) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.count += 1
	c.total += in
	if in < c.min {
		c.min = in
		return
	} else if in > c.max {
		c.max = in
	}
}

func (c *city) getAvg() float64 {
	return c.total / float64(c.count)
}

type mapHandler struct {
	mapping map[string]*city
	mu      sync.RWMutex
}

func (handler *mapHandler) process(name string, in string) {
	c, exist := handler.mapping[name]

	if !exist {
		c = NewCity()
		handler.mu.Lock()
		handler.mapping[name] = c
		handler.mu.Unlock()
	}

	float, err := strconv.ParseFloat(in, 64)
	check(err)

	c.process(float)
}

func (handler *mapHandler) getSortedKeys() []string {
	keys := make([]string, 0, len(handler.mapping))
	for k := range handler.mapping {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func (handler *mapHandler) getCity(name string) *city {
	return handler.mapping[name]
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func ProcessChunk(handler *mapHandler, lines []string) {
	for _, line := range lines {
		line := strings.Split(line, ";")
		city, temp := line[0], line[1]
		handler.process(city, temp)
	}
}
