package main

import "fmt"

func main() {
	result := wrap(hello)
	fmt.Println("result: ", result)

}

func hello() string {
	return "hello there"
}

func wrap(stringer func() string) string {
	return "prefix-" + stringer() + "-suffix"
}

// i want a function that adds a prefix to a string
