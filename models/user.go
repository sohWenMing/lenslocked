package models

import (
	"database/sql"
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type UserToPlainTextPassword struct {
	Email             string
	PlainTextPassword string
}

type User struct {
	ID           int
	Email        string
	PasswordHash string
	*Session
}

type LoggedInUserInfo struct {
	ID    int
	Email string
}

type UserService struct {
	db *sql.DB
	*SessionService
}

func (us *UserService) CreateUser(newUserToCreate UserToPlainTextPassword) (*User, error) {
	preppedInfo := prepUserToPlainTextPassword(newUserToCreate)
	hash, err := generateBcryptHash(preppedInfo.PlainTextPassword)
	if err != nil {
		return nil, err
	}
	//first attempt to generate the hash for the password, hold for storage

	row := us.db.QueryRow(`
		INSERT INTO users (email, password_hash)
		VALUES ($1, $2)
		RETURNING id, email;
	`, preppedInfo.Email, hash)
	returnedUser := User{}

	err = row.Scan(&returnedUser.ID, &returnedUser.Email)
	if err != nil {
		return &User{}, HandlePgError(err)
	}

	session, err := us.SessionService.Create(returnedUser.ID)
	if err != nil {
		return &User{}, errors.New("error occured when trying to create session")
	}
	returnedUser.Session = session

	return &returnedUser, nil
}
func generateBcryptHash(plainTextPassword string) (hash string, err error) {
	hashBytes, err := bcrypt.GenerateFromPassword(
		[]byte(plainTextPassword),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return "", err
	}
	return string(hashBytes), nil
}

func (us *UserService) LoginUser(userToPassword UserToPlainTextPassword) (loggedInUserInfo LoggedInUserInfo, err error) {
	preppedInfo := prepUserToPlainTextPassword(userToPassword)
	row := us.db.QueryRow(`
		SELECT id, email, password_hash 
		FROM  users
		WHERE email=($1);
	`, preppedInfo.Email)
	var returnedUser User
	err = row.Scan(&returnedUser.ID, &returnedUser.Email, &returnedUser.PasswordHash)
	if err != nil {
		return loggedInUserInfo, HandlePgError(err)
	}
	err = bcrypt.CompareHashAndPassword([]byte(returnedUser.PasswordHash), []byte(userToPassword.PlainTextPassword))
	if err != nil {
		return loggedInUserInfo, HandlerBcryptErr(err)
	}
	loggedInUserInfo.ID = returnedUser.ID
	loggedInUserInfo.Email = returnedUser.Email
	return loggedInUserInfo, nil
}

// ##### helpers #####
func prepUserToPlainTextPassword(u UserToPlainTextPassword) UserToPlainTextPassword {
	u.Email = strings.ToLower(u.Email)
	return u
}
