package main

import "fmt" 

func reverseSlice(slice []int){
	for i:= len(slice); i >= len(slice); i--{
		revertedSlice := slice[i - 1]

		return revertedSlice
	}
}

func main(){
	sliceValue := []int{1,2,3,4,5}

	fmt.Println("Inverted Values:", reverseSlice(sliceValue))
}