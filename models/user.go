package models

import (
	"database/sql"
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
}

type LoggedInUserInfo struct {
	ID    int
	Email string
}

type UserService struct {
	db *sql.DB
}

func (us *UserService) CreateUserTableIfNotExist() {
	_, err := us.db.Exec(`
			CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			email TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL
			);
		`)
	if err != nil {
		panic(err)
	}
}

func (us *UserService) CreateUser(newUserToCreate UserToPlainTextPassword) (*User, error) {
	preppedInfo := prepUserToPlainTextPassword(newUserToCreate)
	hashBytes, err := bcrypt.GenerateFromPassword(
		[]byte(preppedInfo.PlainTextPassword),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return nil, err
	}
	hash := string(hashBytes)
	row := us.db.QueryRow(`
		INSERT INTO users (email, password_hash)
		VALUES ($1, $2)
		RETURNING id, email, password_hash;
	`, preppedInfo.Email, hash)
	returnedUser := User{}
	err = row.Scan(&returnedUser.ID, &returnedUser.Email, &returnedUser.PasswordHash)
	if err != nil {
		return &User{}, HandlePgError(err)
	}
	return &returnedUser, nil
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
