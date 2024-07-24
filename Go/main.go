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

// IDEA: Don't use structs and try normal vars and funcs

const KBs = 1024
const MBs = 1024 * KBs

func main() {
	file, err := os.Create("1BRC.prof")
	check(err)
	pprof.StartCPUProfile(file)
	defer pprof.StopCPUProfile()

	countPtr := flag.Int("count", 10, "Number of tests to run and average")
	bufferPtr := flag.Int("buffer", 8, "MB to use for buffer")
	testPtr := flag.Bool("test", false, "Run the smaller input file")
	flag.Parse()

	testCount := *countPtr

	bufferSize := int(*bufferPtr)

	fmt.Println("Running initial pass before starting timing")
	Run1BRC(*testPtr, bufferSize*MBs)
	fmt.Printf("Processing %d tests.\n", testCount)

	start := time.Now()

	for i := 0; i < testCount; i++ {
		Run1BRC(*testPtr, bufferSize*MBs)
	}

	elapsed := time.Since(start)

	average := elapsed.Seconds() / float64(testCount)
	fmt.Printf("Took a total of %s\n", elapsed)
	fmt.Printf("Took an average of %.3fs\n", average)
}

type city struct {
	Min int64
	Max int64
	Sum int64
	Cnt int64
}

func Run1BRC(test bool, bufferSize int) {
	var input *os.File
	var err error

	if test {
		input, err = os.Open("../test_measurements_small.txt")
	} else {
		input, err = os.Open("../measurements.txt")
	}
	check(err)
	defer input.Close()

	var wg sync.WaitGroup

	lineBuffer := make([]byte, bufferSize)
	fileReader := bufio.NewReader(input)

	remainder := make([]byte, 0)
	cityChan := make(chan map[string]city)

	for {
		num, err := fileReader.Read(lineBuffer)
		if num == 0 {
			if err == io.EOF {
				break
			}
			break
		}

		bytesRead := lineBuffer[:num]
		chunk := append(remainder, bytesRead...)

		// ISSUE: Byte splitting like this is taking 13.7%
		splitChunk := bytes.Split(chunk, []byte{'\n'})

		chunkLineCount := len(splitChunk) - 1
		remainder = splitChunk[chunkLineCount]
		splitChunk = splitChunk[:chunkLineCount]

		wg.Add(1)
		go func(splitChunk [][]byte) {
			defer wg.Done()
			ProcessChunk(splitChunk, cityChan)
		}(splitChunk)
	}

	if remainder != nil {
		wg.Add(1)
		go func(splitChunk []byte) {
			defer wg.Done()
			ProcessChunk([][]byte{splitChunk}, cityChan)
		}(remainder)
	}

	go func() {
		wg.Wait()
		close(cityChan)
	}()

	totals := make(map[string]city)
	// fmt.Println("Starting to read from chunk channel")
	for chunk := range cityChan {
		for name, inCity := range chunk {
			c, ok := totals[name]
			if !ok {
				totals[name] = city{
					Min: inCity.Min,
					Max: inCity.Max,
					Sum: inCity.Sum,
					Cnt: inCity.Cnt,
				}
				continue
			}
			c.Min = min(c.Min, inCity.Min)
			c.Max = max(c.Max, inCity.Max)
			c.Sum += inCity.Sum
			c.Cnt += inCity.Cnt
			totals[name] = c
		}
	}

	names := make([]string, 0, len(totals))
	for n := range totals {
		names = append(names, n)
	}

	sort.Strings(names)

	for _, n := range names {
		c := totals[n]

		fmt.Println(fmt.Sprintf("%s=%.1f/%.1f/%.1f", n, float64(c.Min)*0.1, (float64(c.Sum)/float64(c.Cnt))*0.1, float64(c.Max)*0.1))
	}
}

func ProcessChunk(lines [][]byte, cityChan chan map[string]city) {
	chunkMap := make(map[string]city)
	for _, line := range lines {
		lineLen := len(line)
		if lineLen < 2 {
			break
		}

		var semi int
		index := lineLen - 4 // 4 back is the first one that could possibly be a semicolon
		if line[index] == ';' {
			semi = index
		} else if line[index-1] == ';' {
			semi = index - 1
		} else if line[index-2] == ';' {
			semi = index - 2
		}

		name := string(line[:semi])
		byteTemp := append(line[semi+1:lineLen-2], line[lineLen-1])

		// ISSUE: ParseInt is taking 7.1%
		// IDEA:  Manually parse the temp bytes into an int directly, skip strings entirely for temp
		temp, _ := strconv.ParseInt(string(byteTemp), 10, 64)

		c, ok := chunkMap[name]

		if !ok {
			c = city{

				Cnt: 1,
				Max: temp,
				Min: temp,
				Sum: temp,
			}
		} else {
			c.Cnt++
			c.Max = max(c.Max, temp)
			c.Min = min(c.Min, temp)
			c.Sum += temp
		}
		chunkMap[name] = c
	}
	// fmt.Println("Putting chunk map on channel")
	cityChan <- chunkMap
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
