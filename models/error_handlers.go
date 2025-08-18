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

type sqlNoRowsErrEnum int

const (
	NoUserWithEmail sqlNoRowsErrEnum = iota
	NoUserWithUserId
	NoRowsErrorOnRefreshSession
	NewSessionNotReturned
)

func (e sqlNoRowsErrEnum) String() string {
	switch e {
	case NoUserWithEmail:
		return "No user could be found with that email address"
	case NoUserWithUserId:
		return "No user could be found with that user id"
	case NoRowsErrorOnRefreshSession:
		return "no rows were returned when attempting to refresh session"
	case NewSessionNotReturned:
		return "no session was returned after attempting to create session"
	default:
		return "unrecognized error, please check actual error"
	}
}

type sqlNoRowsErrStruct struct {
	enum sqlNoRowsErrEnum
}

func UserNotFoundByUserIdErr() *sqlNoRowsErrStruct {
	return &sqlNoRowsErrStruct{
		NoUserWithUserId,
	}
}
func UserNotFoundByEmailErr() *sqlNoRowsErrStruct {
	return &sqlNoRowsErrStruct{
		NoUserWithEmail,
	}
}
func NoRowsErrorOnRefreshSessionErr() *sqlNoRowsErrStruct {
	return &sqlNoRowsErrStruct{
		NoRowsErrorOnRefreshSession,
	}
}
func NewSessionNotReturnedErr() *sqlNoRowsErrStruct {
	return &sqlNoRowsErrStruct{
		NewSessionNotReturned,
	}
}

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

func HandlePgError(err error, noRowsErr *sqlNoRowsErrStruct) error {
	if errors.Is(err, sql.ErrNoRows) {
		if noRowsErr == nil {
			return MapHandledGenericError(err)
		} else {
			return MapHandledError(err, noRowsErr.enum.String())
		}
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
