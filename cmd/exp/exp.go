package main

import (
	"fmt"
	"path/filepath"
)

func main() {
	pattern := "../../images/gallery-59/*"
	allFiles, err := filepath.Glob(pattern)
	if err != nil {
		fmt.Println(err.Error())
	}
	for i, file := range allFiles {
		fmt.Printf("%d: %s\n", i, file)
		extension := filepath.Ext(file)
		fmt.Println("extension: ", extension)
	}
}
