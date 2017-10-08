package main

import "fmt"
import "time"

// MaxGoroutines -- define the max goroutines number
const MaxGoroutines int = 10

var guard = make(chan struct{}, MaxGoroutines)

func main() {
	for i := 0; i < 30; i++ {
		guard <- struct{}{} // would block if guard channel is already filled
		go worker(i)
	}
}

func worker(i int) {
	time.Sleep(5 * time.Second)
	fmt.Println("doing work on", i)
	<-guard
}
