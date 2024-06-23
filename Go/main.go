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

func NewCity() *city {
	c := city{}

	c.max = -1e99
	c.min = 1e99

	return &c
}

func (c *city) process(in float64) {
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
	mapping map[string]*city
	mu      sync.Mutex
}

func (handler *mapHandler) process(name string, in string) {
	handler.mu.Lock()
	defer handler.mu.Unlock()

	c, exist := handler.mapping[name]
	if !exist {
		c = NewCity()
		handler.mapping[name] = c
	}

	float, err := strconv.ParseFloat(in, 64)
	check(err)

	c.process(float)
}

func (handler *mapHandler) getSortedKeys() []string {
	handler.mu.Lock()
	defer handler.mu.Unlock()

	keys := make([]string, 0, len(handler.mapping))
	for k := range handler.mapping {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func (handler *mapHandler) getCity(name string) *city {
	handler.mu.Lock()
	defer handler.mu.Unlock()

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
	Run1BRC(false)
	pprof.StopCPUProfile()
}

const ChunkSize = 2048
const NumWorkers = 4

func Run1BRC(test bool) {
	var input *os.File
	var err error

	if test {
		input, err = os.Open("../test_measurements.txt")
	} else {
		input, err = os.Open("../measurements.txt")
	}
	check(err)
	defer input.Close()

	lines := make(chan string, 100)
	var wg sync.WaitGroup

	handler := mapHandler{mapping: make(map[string]*city)}

	for i := 0; i < NumWorkers; i++ {
		wg.Add(1)
		go ProcessLine(&handler, &wg, lines)
	}

	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		lines <- scanner.Text()
	}

	close(lines)

	wg.Wait()

	keys := handler.getSortedKeys()

	for k := range keys {
		c := handler.getCity(keys[k])

		fmt.Println(fmt.Sprintf("%s %.1f %.1f %.1f", keys[k], c.min, c.getAvg(), c.max))
	}
}

func ProcessLine(handler *mapHandler, wg *sync.WaitGroup, lines chan string) {
	defer wg.Done()

	for line := range lines {
		line := strings.Split(line, ";")
		city, temp := line[0], line[1]
		handler.process(city, temp)
	}
}
