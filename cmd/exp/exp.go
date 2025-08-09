package main

import (
	"fmt"
	"time"
)

func main() {
	startTime := time.Now()
	combinations := getBitCombinations(8 * 6)
	for i, row := range combinations {
		fmt.Printf("%d: %v\n", i, row)
	}
	timeTaken := time.Since(startTime).Seconds()
	fmt.Printf("Time Taken %v", timeTaken)

	// i want to exhaust all possibilites
}

func getBitCombinations(numBytes int) [][]int {
	if numBytes < 1 {
		return [][]int{}
	}
	return AppendBitRecursive([][]int{
		{0}, {1},
	}, 1, numBytes)
}

func AppendBitRecursive(currentSlices [][]int, curIdx int, lastIdx int) (returnedSlices [][]int) {
	if curIdx == lastIdx {
		return currentSlices
	}
	bitSlice := []int{0, 1}
	workingSlices := [][]int{}
	for _, slice := range currentSlices {

		// here we would isolate each slice, so for example [0], or [1]
		for _, bit := range bitSlice {
			appendedSlice := append(append([]int{}, slice...), bit)
			workingSlices = append(workingSlices, appendedSlice)
			//[0, 0]
		}
	}
	return AppendBitRecursive(workingSlices, curIdx+1, lastIdx)
}
