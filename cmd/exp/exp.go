package main

import (
	"errors"
	"fmt"
)

type Booger struct {
	personName string
}

func (b Booger) Error() string {
	return fmt.Sprintf("this gross booger was from %s", b.personName)
}

func main() {
	boogerErr := &Booger{"Harper"}
	err := fmt.Errorf("This is a test to wrap a booger %w", boogerErr)

	fmt.Println("Is err boogerErr", errors.Is(err, boogerErr))

	var boogerAs *Booger
	if errors.As(err, &boogerAs) {
		fmt.Println(fmt.Sprintf("The error occured because %s made the booger", boogerAs.personName))
	} else {
		fmt.Println("Its a mystery, this thing wasn't actually a booger error")
	}
}
