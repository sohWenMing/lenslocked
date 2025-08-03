package models

import (
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

const (
	genericErrorMsg string = "A problem occured during the operation. Please try again later."
)

type HandledError struct {
	err    error
	errMsg string
}

func (g *HandledError) Error() string {
	return g.errMsg
}

var (
	ErrEmailTaken = errors.New("email has already been used. Please try using another.")
)

func HandlePgError(err error) error {
	if strings.Contains(err.Error(), `ERROR: duplicate key value violates unique constraint "users_email_key" (SQLSTATE 23505)`) {
		return ErrEmailTaken
	} else {
		return &HandledError{err, genericErrorMsg}
	}
}

func HandlerBcryptErr(err error) error {
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return &HandledError{err, "email and password combination do not match"}
	}
	return err
}
