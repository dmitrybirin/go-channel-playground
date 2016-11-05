package main

import (
	"flag"
	"fmt"
	"math"

	"github.com/gosuri/uiprogress"
)

func main() {

	poolCapacity := flag.Int("t", 1, "number of threads")
	barSwitch := flag.Bool("b", true, "show progress bar")
	flag.Parse()

	collection := make([]int, 20)
	var results = make([]float64, 0)

	lenOfNames := len(collection)

	var bar *uiprogress.Bar
	if *barSwitch {
		bar = barInit(lenOfNames)
	}

	resultsChan := make(chan float64, lenOfNames)
	poolChan := make(chan bool, *poolCapacity)

	for i, item := range collection {
		item = i
		go doLongOperation(item, &resultsChan, poolChan)
		poolChan <- true
		if *barSwitch {
			bar.Incr()
		}
	}

	//real-time results retrieving
	for _ = range collection {
		r := <-resultsChan
		results = append(results, r)
	}

	fmt.Println("All Computations Done.")
	fmt.Println("Results are:")
	//after all results printing
	for item := range results {
		fmt.Println(item)
	}
}

func doLongOperation(item int, r *chan float64, p chan bool) {
	defer func() { <-p }()
	result := 0.0
	for i := 0; i < 100000000; i++ {
		result += math.Pi * math.Sin(float64(item))
	}
	*r <- 5.0
}

func barInit(lenOfNames int) *uiprogress.Bar {
	uiprogress.Start()
	bar := uiprogress.AddBar(lenOfNames)
	bar.AppendCompleted()
	bar.PrependElapsed()
	return bar
}
