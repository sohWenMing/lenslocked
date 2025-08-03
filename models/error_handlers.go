package models

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrEmailTaken     = errors.New("email has already been used. Please try using another.")
	ErrGenericPGError = errors.New("A generic error has occured, please try again later")
)

func HandlePgError(err error) error {
	fmt.Println("error: ", err)
	if strings.Contains(err.Error(), `ERROR: duplicate key value violates unique constraint "users_email_key" (SQLSTATE 23505)`) {
		return ErrEmailTaken
	} else {
		return fmt.Errorf("%s %w", err.Error(), ErrGenericPGError)
	}
}
