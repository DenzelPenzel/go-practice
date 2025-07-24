package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"runtime/trace"
	"time"

	"github.com/DenzelPenzel/go-leetcode/interview/billion-rows/sol1"
)

var name = flag.String("name", "", "path to the file")
var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
var memprofile = flag.String("memprofile", "", "write memory profile to `file`")
var executionprofile = flag.String("execprofile", "", "write trace execution to `file`")

func main() {
	flag.Parse()

	start := time.Now()

	if *executionprofile != "" {
		f, err := os.Create("./prof/" + *executionprofile)
		if err != nil {
			log.Fatalf("Failed to trace app prof: %v", err)
		}
		defer f.Close()
		err = trace.Start(f)
		if err != nil {
			log.Fatalf("Failed to start trace")
		}
		defer trace.Stop()
	}

	if *cpuprofile != "" {
		f, err := os.Create("./prof/" + *cpuprofile)
		if err != nil {
			log.Fatalf("Failed to create CPU prof %v", err)
		}
		defer f.Close()
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatalf("Failed to start CPU prof %v", err)
		}
		defer pprof.StopCPUProfile()
	}

	if *name == "" {
		log.Fatalf("Filename param is missing")
	}

	sol1.Run("./data/" + *name)

	fmt.Println(time.Now().Sub(start))

	if *memprofile != "" {
		f, err := os.Create("./prof/" + *memprofile)
		if err != nil {
			log.Fatalf("Failed to create mem prof %v", err)
		}
		defer f.Close()
		runtime.GC()
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatalf("failed to write mem prof %v", err)
		}
	}
}
