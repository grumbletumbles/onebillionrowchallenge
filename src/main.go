package main

import (
	"runtime"
	"sync"
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

func worker(batchChan <-chan []string, resultChan chan<- map[string]Data, wg *sync.WaitGroup) {
	defer wg.Done()

}

func main() {
	numWorkers := runtime.NumCPU()
	batchChan := make(chan []string, batchSize)
	resultsChan := make(chan map[string]Data, numWorkers)

	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(batchChan, resultsChan, &wg)
	}

}
