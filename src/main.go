package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	batchSize = 1000000
)

type Data struct {
	Min   float64
	Max   float64
	Sum   float64
	Count int32
}

func (d *Data) Update(value float64) {
	d.Min = min(d.Min, value)
	d.Max = max(d.Max, value)
	d.Sum += value
	d.Count++
}

func NewData(value float64) *Data {
	return &Data{
		Min:   value,
		Max:   value,
		Sum:   value,
		Count: 1,
	}
}

func (d *Data) Add(other *Data) {
	d.Min = min(d.Min, other.Min)
	d.Max = max(d.Max, other.Max)
	d.Sum += other.Sum
	d.Count += other.Count
}

func worker(batchChan <-chan []string, resultChan chan<- map[string]*Data, wg *sync.WaitGroup) {
	defer wg.Done()

	stats := make(map[string]*Data)
	for batch := range batchChan {
		for _, line := range batch {
			x := strings.Split(line, ";")
			city := x[0]
			t, err := strconv.ParseFloat(x[1], 64)
			if err != nil {
				log.Fatal(err)
			}
			if _, ok := stats[city]; ok {
				stats[city].Update(t)
			} else {
				stats[city] = NewData(t)
			}
		}
	}

	resultChan <- stats
}

func main() {
	start := time.Now()

	numWorkers := runtime.NumCPU()
	batchChan := make(chan []string, batchSize)
	resultsChan := make(chan map[string]*Data, numWorkers)

	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(batchChan, resultsChan, &wg)
	}

	go func() {
		file, err := os.Open("data/measurements.txt")
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		var batch []string

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			batch = append(batch, scanner.Text())
			if len(batch) >= batchSize {
				batchChan <- batch
				batch = nil
			}
		}

		if len(batch) >= 0 {
			batchChan <- batch
		}
		close(batchChan)
	}()

	wg.Wait()
	close(resultsChan)

	final := make(map[string]*Data)
	for result := range resultsChan {
		for station, stats := range result {
			if _, ok := final[station]; ok {
				final[station].Add(stats)
			} else {
				final[station] = stats
			}
		}
	}

	for c, t := range final {
		fmt.Printf("%s -- min: %.2f max: %.2f avg: %.2f\n", c, t.Min, t.Max, t.Sum/float64(t.Count))
	}

	fmt.Printf("total time: %s\n", time.Since(start))
}
