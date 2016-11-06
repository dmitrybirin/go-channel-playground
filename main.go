package main

import (
	"flag"
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/gosuri/uiprogress"
)

func main() {

	var poolCapacity int
	var collectionLen int
	var barSwitch bool
	var results = make([]float64, 0)

	flag.IntVar(&poolCapacity, "t", 1, "number of threads")
	flag.IntVar(&collectionLen, "l", 10, "length of initial collection")
	flag.BoolVar(&barSwitch, "p", false, "show progress bar")
	flag.Parse()

	var bar *uiprogress.Bar
	if barSwitch {
		bar = barInit(collectionLen)
	}

	randSource := time.Now().UnixNano()
	r := rand.New(rand.NewSource(randSource))

	collection := make([]int, collectionLen)

	resultsChan := make(chan float64, collectionLen)
	poolChan := make(chan bool, poolCapacity)

	for _, item := range collection {
		item = r.Int()
		go doLongOperation(item, resultsChan, poolChan)
		poolChan <- true
		if barSwitch {
			bar.Incr()
		}
	}

	//real-time results retrieving
	for _ = range collection {
		results = append(results, <-resultsChan)
	}

	fmt.Println("All Computations Done.")
	fmt.Println("Results are:")
	//after all results printing
	for _, item := range results {
		fmt.Println(item)
	}
}

func doLongOperation(item int, r chan float64, p chan bool) {
	defer func() { <-p }()
	result := 0.0
	for i := 0; i < 100000000; i++ {
		result += math.Pi * math.Sin(float64(item))
	}
	r <- result
}

func barInit(lenOfNames int) *uiprogress.Bar {
	uiprogress.Start()
	bar := uiprogress.AddBar(lenOfNames)
	bar.AppendCompleted()
	bar.PrependElapsed()
	return bar
}
