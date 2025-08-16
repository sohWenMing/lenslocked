package models

import (
	"database/sql"
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

const (
	genericErrorMsg    string = "A problem occured during the operation. Please try again later."
	emailTakenErrorMsg string = "email has already been used. Please try using another."
)

type HandledError struct {
	err    error
	errMsg string
}

func (g *HandledError) Error() string {
	return g.errMsg
}

func MapHandledError(err error, errMsg string) *HandledError {
	return &HandledError{
		err, errMsg,
	}
}
func MapHandledGenericError(err error) *HandledError {
	return &HandledError{
		err, genericErrorMsg,
	}
}

/*
Handles case where no rows are returned from query, but by design it's not an error.
Returns true if no rows are found
*/
func CheckIsNoRowsErr(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}

func HandlePgError(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return MapHandledError(err, "No user could be found with that email address")
	}
	if strings.Contains(err.Error(), `ERROR: duplicate key value violates unique constraint "users_email_key" (SQLSTATE 23505)`) {
		return MapHandledError(err, emailTakenErrorMsg)
	} else {
		return MapHandledGenericError(err)
	}
}

func HandlerBcryptErr(err error) error {
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return MapHandledError(err, "email and password combination do not match")
	}
	return MapHandledGenericError(err)
}
