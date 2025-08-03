package main

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	for i, arg := range os.Args {
		fmt.Println("arg :", i, arg)
	}

	if len(os.Args) <= 1 {
		panic("number of arguments passed in must be more than on equal to 1")
	}

	operation := os.Args[1]
	if checkStringsCompare("compare", operation) {
		compareFunc()
		return
	}
	if checkStringsCompare("hash", operation) {
		hash, err := hashFunc()
		if err != nil {
			panic(err)
		}
		askUserCheckPassword(hash)
		return
	}
	fmt.Printf("operation %s is not supported, please enter compare or hash as first argument", operation)
	os.Exit(0)
}

func checkStringsCompare(checkValue, inputValue string) bool {
	return strings.ToUpper(strings.TrimSpace(checkValue)) == strings.ToUpper(strings.TrimSpace(inputValue))
}

func compareFunc() {
	fmt.Println("compare func called")
}
func hashFunc() (hash string, err error) {
	if len(os.Args) <= 2 {
		fmt.Println("password needs to passed in as second argument when hashing")
		return
	}
	password := os.Args[2]
	fmt.Println("password passed in: ", password)
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	hash = string(hashBytes)
	fmt.Println("hash generated: ", hash)
	return hash, nil
}

func askUserCheckPassword(hash string) {
	var password string
	fmt.Println("Enter password to check")
	fmt.Scanln(&password)
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		fmt.Println("password not valid")
		return
	} else {
		fmt.Println("password is valid")
	}
}
