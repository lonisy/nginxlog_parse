package main

import (
	"fmt"
	"runtime"
	"time"
)

func main() {
	fmt.Println(runtime.NumCPU())
	runtime.GOMAXPROCS(runtime.NumCPU())
	for i := 0; i < 100000; i++ {
		go func(i int) {
			for {
				fmt.Println(i)
			}
		}(i)
	}
	time.Sleep(time.Minute)
}
