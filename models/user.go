package models

import (
	"database/sql"

	"golang.org/x/crypto/bcrypt"
)

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

func (us *UserService) CreateUser(email, password string) (*User, error) {
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	hash := string(hashBytes)
	row := us.db.QueryRow(`
		INSERT INTO users (email, password_hash)
		VALUES ($1, $2)
		RETURNING id, email, password_hash;
	`, email, hash)
	returnedUser := User{}
	err = row.Scan(&returnedUser.ID, &returnedUser.Email, &returnedUser.PasswordHash)
	if err != nil {
		return &User{}, HandlePgError(err)
	}
	return &returnedUser, nil
}

func (us *UserService) LoginUser(email, password string) (loggedInUserInfo LoggedInUserInfo, err error) {

	row := us.db.QueryRow(`
		SELECT id, email, password_hash 
		FROM  users
		WHERE email=($1);
	`, email)
	var returnedUser User
	err = row.Scan(&returnedUser.ID, &returnedUser.Email, &returnedUser.PasswordHash)
	if err != nil {
		return loggedInUserInfo, HandlePgError(err)
	}
	err = bcrypt.CompareHashAndPassword([]byte(returnedUser.PasswordHash), []byte(password))
	if err != nil {
		return loggedInUserInfo, HandlerBcryptErr(err)
	}
	loggedInUserInfo.ID = returnedUser.ID
	loggedInUserInfo.Email = returnedUser.Email
	return loggedInUserInfo, nil
}
