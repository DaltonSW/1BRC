package main

import (
	"bufio"
	"flag"
	"fmt"
	"sync"
	"time"

	"os"
	"runtime/pprof"
	"sort"
	"strconv"
	"strings"
)

// NOTE: 1BRC Requirements
//	Each row: <string: station name>;<double: measurement with EXACTLY one fractional digit>
//	Must calculate min, max, and average for each city
//	Must output them to stdout, alphabetically sorted by city name, in the format "name=min/mean/max"

type city struct {
	count int
	total int
	min   int
	max   int
	mu    sync.RWMutex
}

func NewCity() *city {
	c := city{}

	c.max = -99999
	c.min = 99999

	return &c
}

func (c *city) process(in int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.count += 1
	c.total += in
	c.min = min(c.min, in)
	c.max = max(c.max, in)
}

func (c *city) getAvg() int {
	return c.total / c.count
}

type mapHandler struct {
	mapping map[string]*city
	mu      sync.RWMutex
}

func (handler *mapHandler) process(name string, numIn string) {
	c, exist := handler.mapping[name]

	if !exist {
		c = NewCity()
		handler.mu.Lock()
		handler.mapping[name] = c
		handler.mu.Unlock()
	}

	numLen := len(numIn)
	tempAsInt, _ := strconv.Atoi(numIn[:numLen-2] + string(numIn[numLen-1]))
	c.process(tempAsInt)
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

const KBs = 1024
const MBs = 1024 * KBs
const TestCount = 1

const DecimalMult = 0.1

func startProfiling() {
	file, err := os.Create("1BRC.prof")
	if err != nil {
		fmt.Println("Couldn't open profiling file")
		os.Exit(1)
	}
	pprof.StartCPUProfile(file)
}

func main() {
	profPtr := flag.Bool("prof", false, "Profile the program")

	flag.Parse()
	if *profPtr {
		startProfiling()
		defer pprof.StopCPUProfile()
	}

	start := time.Now()
	for i := 0; i < TestCount; i++ {
		Run1BRC(false, 1*MBs)
	}
	elapsed := time.Since(start)
	average := elapsed / TestCount
	fmt.Printf("Processing %d tests.\n", TestCount)
	fmt.Printf("Took a total of %s\n", elapsed)
	fmt.Printf("Took an average of %s\n", average)
}

func Run1BRC(test bool, bufferSize int) {
	var input *os.File
	var err error

	if test {
		input, err = os.Open("../test_measurements.txt")
	} else {
		input, err = os.Open("../measurements.txt")
	}
	if err != nil {
		fmt.Println("Couldn't open measurements file!")
		os.Exit(1)
	}
	defer input.Close()

	lineBuffer := make([]byte, bufferSize)
	fileReader := bufio.NewReader(input)

	var wg sync.WaitGroup
	handler := mapHandler{mapping: make(map[string]*city)}
	remainder := ""

	for {
		num, _ := fileReader.Read(lineBuffer)
		if num == 0 {
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
		fmt.Println(fmt.Sprintf("%s=%.1f/%.1f/%.1f", keys[k], float64(c.min)*DecimalMult, float64(c.getAvg())*DecimalMult, float64(c.max)*DecimalMult))
	}
}

func ProcessChunk(handler *mapHandler, lines []string) {
	for _, line := range lines {
		line := strings.Split(line, ";")
		// if len(line) == 1 {
		// 	log.Error("Line couldn't be parsed!", "line", line)
		//
		// }
		city, temp := line[0], line[1]
		handler.process(city, temp)
	}
}
