package main

import (
	"fmt"
)

func readStrings(strings ...string) {
	for i, text := range strings {
		fmt.Println(fmt.Sprintf("index %d: %s", i, text))
	}
}
func main() {
	strings := []string{
		"one", "two",
	}
	readStrings(strings...)
}
