package main

import (
	"fmt"
	"runtime/debug"
	"time"
)

func periodicFree(d time.Duration) {
	tick := time.Tick(d)
	for _ = range tick {
		debug.FreeOSMemory()
	}
}

func main() {
	fmt.Println(111)
}
