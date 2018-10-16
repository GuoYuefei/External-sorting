package main

import (
	"fmt"
	"time"
)

func main() {
	for i := 0; i < 1000; i++ {
		func(ii int) {
			go printHelloWorld(i)
		}(i)
	}
	time.Sleep(200 * time.Millisecond)
}

func printHelloWorld(i int) {
	for {
		fmt.Printf("hello world from goroutine %d !\n", i)
	}
}
