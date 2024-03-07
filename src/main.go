package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
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

func main() {
	start := time.Now()
	f, err := os.Open("data/measurements.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	temp := make(map[string]*Data)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		str := scanner.Text()
		x := strings.Split(str, ";")
		city := x[0]
		t, err := strconv.ParseFloat(x[1], 64)
		if err != nil {
			log.Fatal(err)
		}
		if _, ok := temp[city]; ok {
			temp[city].Update(t)
		} else {
			temp[city] = NewData(t)
		}
	}

	//for c, t := range temp {
	//	fmt.Printf("%s -- min: %.2f max: %.2f avg: %.2f\n", c, t.Min, t.Max, t.Sum/float64(t.Count))
	//}
	fmt.Printf("total time: %s\n", time.Since(start))
}
