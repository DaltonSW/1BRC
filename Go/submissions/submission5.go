package main

import (
	"bufio"
	"fmt"
	"io"
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
	handler.mu.RLock()
	c, exist := handler.mapping[name]
	handler.mu.RUnlock()

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

const KBs = 1024
const MBs = 1024 * KBs

func main() {
	file, err := os.Create("1BRC.prof")
	check(err)
	pprof.StartCPUProfile(file)
	defer pprof.StopCPUProfile()
	Run1BRC(false, 8*MBs)
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

		fmt.Println(fmt.Sprintf("%s %.1f %.1f %.1f", keys[k], c.min, c.getAvg(), c.max))
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
