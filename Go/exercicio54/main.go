package main

import (
	"fmt"
	"sync"
	"time"
)

func SayMayName(s string) {

	for i := 0; i < 5; i++ {
		time.Sleep(100 * time.Millisecond)
		fmt.Println(s)
	}
}

func main() {
	startNow := time.Now()
	var (
		mayName string = "Ricardo"
		wg      sync.WaitGroup
	)
	fmt.Println("With gourootines---------------------------------------- ")
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			SayMayName(mayName)
			defer wg.Done()
		}()
	}
	wg.Wait()

	fmt.Println("Operation time: ", time.Since(startNow))
}
