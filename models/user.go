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
		returning id, email, password_hash;
	`, email, hash)
	returnedUser := User{}
	err = row.Scan(&returnedUser.ID, &returnedUser.Email, &returnedUser.PasswordHash)
	if err != nil {
		return &User{}, HandlePgError(err)
	}
	return &returnedUser, nil
}
