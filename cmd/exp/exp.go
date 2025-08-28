package main

import (
	"context"
	"errors"
	"fmt"
)

type contextColor string

func main() {
	ctx := context.WithValue(context.Background(), contextColor("favourite-color"), "blue")
	value := ctx.Value(contextColor("favourite-color"))
	fmt.Println("value:", value)
	valueString, ok := value.(string)
	if !ok {
		panic(errors.New("value from favourite color is not a string"))
	}
	fmt.Println("valueString", valueString)
}
