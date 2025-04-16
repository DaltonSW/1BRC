package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"sync"
	"time"
	"unicode"

	"os"
	"runtime/pprof"
	"sort"
)

const KBs = 1024
const MBs = 1024 * KBs

func main() {
	// Flags
	countPtr := flag.Int("count", 10, "Number of tests to run and average")
	bufferPtr := flag.Int("buffer", 16, "MB to use for buffer")
	testPtr := flag.Bool("test", false, "Run the smaller input file")
	flag.Parse()
	testCount := *countPtr
	bufferSize := int(*bufferPtr)

	var filepath string
	if *testPtr {
		filepath = "../test_measurements.txt"
	} else {
		filepath = "../measurements.txt"
	}

	// Initial run
	fmt.Println("Running initial pass before starting timing")
	Run1BRC(bufferSize*MBs, filepath)

	// Start profiling
	file, err := os.Create("1BRC.prof")
	check(err)
	pprof.StartCPUProfile(file)
	defer pprof.StopCPUProfile()

	// Test running
	start := time.Now()
	fmt.Printf("Processing %d tests.\n", testCount)
	for i := 0; i < testCount; i++ {
		Run1BRC(bufferSize*MBs, filepath)
	}
	elapsed := time.Since(start)

	// Print info
	average := elapsed.Seconds() / float64(testCount)
	fmt.Printf("Took a total of %s\n", elapsed)
	fmt.Printf("Took an average of %.3fs\n", average)
}

type city struct {
	Min int32
	Max int32
	Sum int32
	Cnt int
}

func Run1BRC(bufferSize int, filepath string) {
	// var declarations
	var input *os.File
	var err error
	var wg sync.WaitGroup

	lineBuffer := make([]byte, bufferSize)
	remainder := make([]byte, 0, bufferSize)
	cityChan := make(chan map[string]*city, 10)

	// Open measurements
	input, err = os.Open(filepath)
	check(err)
	defer input.Close()
	fileReader := bufio.NewReader(input)

	// Loop reading
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

		splitChunk := bytes.Split(chunk, []byte{'\n'}) // PERF: 12.1% -- bytes.Split

		chunkLineCount := len(splitChunk) - 1
		remainder = splitChunk[chunkLineCount]
		splitChunk = splitChunk[:chunkLineCount]

		wg.Add(1)
		go func(splitChunk [][]byte) {
			ProcessChunk(splitChunk, cityChan)
			wg.Done()
		}(splitChunk)
	}

	if remainder != nil {
		wg.Add(1)
		go func(splitChunk []byte) {
			ProcessChunk([][]byte{splitChunk}, cityChan)
			wg.Done()
		}(remainder)
	}

	go func() {
		wg.Wait()
		close(cityChan)
	}()

	totals := make(map[string]*city)
	for chunk := range cityChan {
		for name, inCity := range chunk {
			c, ok := totals[name]
			if !ok {
				totals[name] = &city{
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

func ProcessChunk(lines [][]byte, cityChan chan map[string]*city) {
	chunkMap := make(map[string]*city)
	for _, line := range lines {
		lineLen := len(line)
		if lineLen < 2 {
			break
		}

		var semi int
		var ones, tens, hundreds, temp int32
		index := lineLen - 1

		ones = int32(line[index] - '0')
		index -= 2 // Skip the period

		tens = int32(line[index] - '0')
		index--

		if unicode.IsDigit(rune(line[index])) {
			hundreds = int32(line[index] - '0')
			index--
		}

		temp = hundreds*100 + tens*10 + ones

		if bytes.Equal([]byte{line[index]}, []byte{'-'}) {
			temp = -temp
			index--
		}

		semi = index

		// PERF: 20.8% -- slicebytetostring
		name := string(line[:semi])

		// PERF: 30.8% -- mapaccess2_faststr
		c, ok := chunkMap[name]
		if !ok {
			chunkMap[name] = &city{

				Cnt: 1,
				Max: temp,
				Min: temp,
				Sum: temp,
			}
			continue
		}

		c.Cnt++
		c.Max = max(c.Max, temp)
		c.Min = min(c.Min, temp)
		c.Sum += temp
	}
	cityChan <- chunkMap
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
