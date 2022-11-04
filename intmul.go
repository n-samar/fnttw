package main

import (
	"flag"
	"fmt"
	"os"
	"runtime/pprof"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	logical_cores := flag.Int("logical_cores", 0, "logical_cores")
	elems_per_core := flag.Int("elems_per_core", 0, "elems_per_core")
	iters := flag.Int("iters", 100, "iterations")

	flag.Parse()

	a := make([]uint64, (*logical_cores)*(*elems_per_core))

	// init
	for i := 0; i < *logical_cores; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			start := i * (*elems_per_core)
			end := (i + 1) * (*elems_per_core)
			for j := start; j < end; j++ {
				a[j] = uint64(j)
			}
		}(i)
	}

	wg.Wait()

	profile_file, _ := os.Create("./intmul.pprof")
	pprof.StartCPUProfile(profile_file)

	for i := 0; i < *logical_cores; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for iter := 0; iter < *iters; iter++ {
				start := i * (*elems_per_core)
				end := (i + 1) * (*elems_per_core)
				for j := start; j < end; j++ {
					a[j] = a[j] * uint64(j)
				}
			}
		}(i)
	}

	wg.Wait()
	pprof.StopCPUProfile()
	fmt.Println("Total multiplies: ", (*logical_cores)*(*elems_per_core)*(*iters))
}
