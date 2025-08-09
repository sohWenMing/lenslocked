package models

import (
	"database/sql"
	"errors"
	"fmt"
	"net/mail"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type UserToPlainTextPassword struct {
	Email             string
	PlainTextPassword string
}

type User struct {
	ID int
	*Session
}
type internalUserStruct struct {
	ID           int
	Email        string
	PasswordHash string
	*Session
}

type UserService struct {
	db *sql.DB
	*SessionService
}

func (us *UserService) CreateUser(newUserToCreate UserToPlainTextPassword) (*User, error) {
	if !isValidEmail(newUserToCreate.Email) {
		return nil, errors.New("Email is not valid")
	}
	if !isValidPassword(newUserToCreate.PlainTextPassword) {
		return nil, errors.New("Password is not valid")
	}
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
	internalUser := internalUserStruct{}

	err = row.Scan(&internalUser.ID, &internalUser.Email)
	if err != nil {
		return &User{}, HandlePgError(err)
	}
	fmt.Println("internal user: ", internalUser)

	session, err := us.SessionService.Create(internalUser.ID)
	if err != nil {
		return &User{}, errors.New("error occured when trying to create session")
	}
	internalUser.Session = session
	returnedUser := &User{
		internalUser.ID, internalUser.Session,
	}
	fmt.Println("returned user: ", returnedUser)

	return returnedUser, nil
}

func isValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return false
	}
	return true
}

func isValidPassword(password string) bool {
	if password == "" {
		return false
	}
	if len(password) < 6 {
		return false
	}
	return true
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

func (us *UserService) LoginUser(userToPassword UserToPlainTextPassword) (user *User, err error) {
	preppedInfo := prepUserToPlainTextPassword(userToPassword)
	row := us.db.QueryRow(`
		SELECT id, email, password_hash 
		FROM  users
		WHERE email=($1);
	`, preppedInfo.Email)
	var internalUser internalUserStruct
	err = row.Scan(&internalUser.ID, &internalUser.Email, &internalUser.PasswordHash)
	if err != nil {
		return nil, HandlePgError(err)
	}
	err = bcrypt.CompareHashAndPassword([]byte(internalUser.PasswordHash), []byte(userToPassword.PlainTextPassword))
	if err != nil {
		return nil, HandlerBcryptErr(err)
	}
	session, err := us.SessionService.ClearPrevSessionsAndCreateNewSessionByUserId(internalUser.ID)
	if err != nil {
		handlerError := HandlePgError(err)
		//TODO change Print to log function
		fmt.Println(err)
		return nil, errors.New(handlerError.Error())
	}
	internalUser.Session = session
	return &User{internalUser.ID, internalUser.Session}, nil
}

func (us *UserService) DeleteUserAndSession(userId int) (err error) {
	err = us.SessionService.DeleteAllSessionsTokensByUserId(userId)
	if err != nil {
		return err
	}
	_, err = us.db.Exec(`
	DELETE from users
	WHERE id=($1);`, userId)
	return err
}

// ##### helpers #####
func prepUserToPlainTextPassword(u UserToPlainTextPassword) UserToPlainTextPassword {
	u.Email = strings.ToLower(u.Email)
	return u
}
