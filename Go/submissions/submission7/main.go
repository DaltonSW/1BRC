package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"strconv"
	"sync"
	"time"

	"os"
	"runtime/pprof"
	"sort"
)

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

	fmt.Printf("Processing %d tests.\n", testCount)
	start := time.Now()
	for i := 0; i < testCount; i++ {
		Run1BRC(false, 8*MBs)
	}
	elapsed := time.Since(start)
	average := elapsed.Seconds() / float64(testCount)
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
	remainder := make([]byte, 0)

	for {
		num, err := fileReader.Read(lineBuffer)
		if num == 0 {
			if err == io.EOF {
				break
			}
			check(err)
			break
		}

		bytesRead := lineBuffer[:num]
		chunk := append(remainder, bytesRead...)
		splitChunk := bytes.Split(chunk, []byte{'\n'})

		chunkLineCount := len(splitChunk) - 1
		remainder = splitChunk[chunkLineCount]
		splitChunk = splitChunk[:chunkLineCount]

		wg.Add(1)
		go func(chunk [][]byte) {
			defer wg.Done()
			ProcessChunk(&handler, chunk)
		}(splitChunk)
	}

	if remainder != nil {
		ProcessChunk(&handler, [][]byte{remainder})
	}

	wg.Wait()

	keys := handler.getSortedKeys()

	for k := range keys {
		c := handler.mapping[keys[k]]

		fmt.Println(fmt.Sprintf("%s=%.1f/%.1f/%.1f", keys[k], float64(c.min)*0.1, c.getAvg()*0.1, float64(c.max)*0.1))
	}
}

func ProcessChunk(handler *mapHandler, lines [][]byte) {
	for _, line := range lines {
		lineLen := len(line)
		if lineLen < 2 {
			break
		}
		// NOTE: Manually locate semicolon instead of using bytes.Split or strings.Split
		var semi int
		index := lineLen - 4 // 4 back is the first one that could possibly be a semicolon
		if line[index] == ';' {
			semi = index
		} else if line[index-1] == ';' {
			semi = index - 1
		} else if line[index-2] == ';' {
			semi = index - 2
		}
		city := line[:semi]
		temp := line[semi+1:]
		handler.process(string(city), string(temp))
	}
}

type city struct {
	count int
	total int32
	min   int32
	max   int32
}

func NewCity() *city {
	c := &city{}

	c.max = -9999
	c.min = 9999

	return c
}

// NOTE: Did some benchmarking, using min() and max() are way faster than manual comparison
func (c *city) process(in int32) {
	c.count += 1
	c.total += in
	c.min = min(c.min, in)
	c.max = max(c.max, in)
}

func (c *city) getAvg() float64 {
	return float64(c.total) / float64(c.count)
}

type mapHandler struct {
	mapping map[string]*city
	mu      sync.RWMutex
}

func (handler *mapHandler) process(name string, inTemp string) {
	c, exist := handler.mapping[name]

	if !exist {
		c = NewCity()
		handler.mu.Lock()
		handler.mapping[name] = c
		handler.mu.Unlock()
	}

	float, err := strconv.ParseFloat(inTemp, 64)
	check(err)

	c.process(int32(float * 10))
}

func (handler *mapHandler) getSortedKeys() []string {
	keys := make([]string, 0, len(handler.mapping))
	for k := range handler.mapping {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
