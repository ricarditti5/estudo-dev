package main

import (
	"fmt"
	"runtime"
)

func printMemory() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	fmt.Printf("\nUsing memory %v", m.Alloc/1024/1024)
}

func main() {

	printMemory()
	memoryAloc := make([]int, 100*1024*1024)
	printMemory()
	//memoryAloc = nil
	//runtime.GC()
	fmt.Println("\n", len(memoryAloc))
}
