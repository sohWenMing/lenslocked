package main

import (
	"errors"
	"fmt"
	"log"
)

func Connect() error {
	//try to connect
	// pretend we got an error
	return errors.New("connection failed")
}

func CreateUser() error {
	err := Connect()
	if err != nil {
		return fmt.Errorf("create user function %w", err)
	}
	return nil
}

func CreateOrg() error {
	err := CreateUser()
	if err != nil {
		return fmt.Errorf("create org function:  %w", err)
	}
	return nil
}

func main() {
	err := CreateUser()
	if err != nil {
		log.Println(err)
	}
	err = CreateOrg()
	if err != nil {
		log.Println(err)
	}
}
