package models

import (
	"database/sql"
	"errors"
	"fmt"
	"net/mail"
	"strings"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

type UserEmailToPlainTextPassword struct {
	Email             string
	PlainTextPassword string
}

type UserIdToSession struct {
	ID int
	*Session
}

type UserInfo struct {
	ID    int
	Email string
}

func (userInfo *UserInfo) String() string {
	return fmt.Sprintf("UserId: %d Email: %s", userInfo.ID, userInfo.Email)
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

func (us *UserService) CreateUser(newUserToCreate UserEmailToPlainTextPassword) (*UserIdToSession, error) {
	err := validateEmailAndPassword(newUserToCreate.Email, newUserToCreate.PlainTextPassword)
	if err != nil {
		return nil, err
	}
	preppedInfo := setEmailLowerCaseInUserToPlainTextPassword(newUserToCreate)
	hash, err := GenerateBcryptHash(preppedInfo.PlainTextPassword)
	if err != nil {
		return nil, err
	}
	//first attempt to generate the hash for the password, hold for storage in the database operation below
	row := us.db.QueryRow(`
		INSERT INTO users (email, password_hash)
		VALUES ($1, $2)
		RETURNING id, email;
	`, preppedInfo.Email, hash)
	internalUser := internalUserStruct{}

	err = row.Scan(&internalUser.ID, &internalUser.Email)
	if err != nil {
		return &UserIdToSession{}, HandlePgError(err, UserNotFoundByEmailErr())
	}

	session, err := us.SessionService.CreateSession(internalUser.ID)
	if err != nil {
		return &UserIdToSession{}, errors.New("error occured when trying to create session")
	}
	internalUser.Session = session
	returnedUser := mapInternalUserToReturnedUser(internalUser)
	return returnedUser, nil
}

func (us *UserService) GetUserByEmail(email string) (userIdToEmail UserInfo, err error) {
	row := us.db.QueryRow(`
		SELECT id, email 
		  FROM users
		 WHERE users.email = ($1);
	`, email)
	var uIdToEmail UserInfo
	err = row.Scan(&uIdToEmail.ID, &uIdToEmail.Email)
	if err != nil {
		return UserInfo{}, HandlePgError(err, UserNotFoundByUserIdErr())
	}
	return uIdToEmail, nil

}

func (us *UserService) GetUserById(userId int) (userIdToEmail UserInfo, err error) {
	row := us.db.QueryRow(`
		SELECT id, email 
		  FROM users
		 WHERE users.id = ($1);
	`, userId)
	var uIdToEmail UserInfo
	err = row.Scan(&uIdToEmail.ID, &uIdToEmail.Email)
	if err != nil {
		return UserInfo{}, HandlePgError(err, UserNotFoundByUserIdErr())
	}
	return uIdToEmail, nil
}

func (us *UserService) LoginUser(userToPassword UserEmailToPlainTextPassword) (user *UserIdToSession, err error) {
	err = validateEmailAndPassword(userToPassword.Email, userToPassword.PlainTextPassword)
	if err != nil {
		return nil, err
	}
	preppedInfo := setEmailLowerCaseInUserToPlainTextPassword(userToPassword)
	row := us.db.QueryRow(`
		SELECT id, email, password_hash 
		FROM  users
		WHERE email=($1);
	`, preppedInfo.Email)
	var internalUser internalUserStruct
	err = row.Scan(&internalUser.ID, &internalUser.Email, &internalUser.PasswordHash)
	if err != nil {
		return nil, HandlePgError(err, UserNotFoundByEmailErr())
	}
	err = bcrypt.CompareHashAndPassword([]byte(internalUser.PasswordHash), []byte(userToPassword.PlainTextPassword))
	if err != nil {
		return nil, HandlerBcryptErr(err)
	}
	session, err := us.SessionService.ExpirePreviousSessionsAndCreateNewSessionByUserId(internalUser.ID)
	if err != nil {
		handlerError := HandlePgError(err, NewSessionNotReturnedErr())
		//TODO change Print to log function
		fmt.Println(err)
		return nil, errors.New(handlerError.Error())
	}
	internalUser.Session = session
	return mapInternalUserToReturnedUser(internalUser), nil
}

func (us *UserService) LogoutUserByToken(token string) error {
	//TODO - implement method
	return nil
}

func (us *UserService) LogoutUser(userId int) (err error) {
	err = us.SessionService.ExpireSessionsTokensByUserId(userId)
	if err != nil {
		return err
	}
	return nil
}

func (us *UserService) GetUserCountById(userId int) (count int, err error) {

	row := us.db.QueryRow(`
		SELECT COUNT(*)
		FROM users
		WHERE id=($1)
	`, userId)
	err = row.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (us *UserService) DeleteUserAndSession(userId int) (err error) {
	_, err = us.db.Exec(`
	DELETE from users
	WHERE id=($1);`, userId)
	return err
}

// preps the UserToPlainTextPassword struct by lowercasing the email to ensure consistency
func setEmailLowerCaseInUserToPlainTextPassword(u UserEmailToPlainTextPassword) UserEmailToPlainTextPassword {
	u.Email = strings.ToLower(u.Email)
	return u
}

func validateEmailAndPassword(email, password string) error {
	if !isValidEmail(email) {
		return errors.New("email is not valid")
	}
	if !isValidPassword(password) {
		return errors.New("password is not valid")
	}
	return nil
}

func isValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
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

func mapInternalUserToReturnedUser(internalUser internalUserStruct) *UserIdToSession {
	returnedUser := &UserIdToSession{
		internalUser.ID, internalUser.Session,
	}
	return returnedUser
}

// ##### helpers #####
func CleanUpCreatedUserIds(createdUserIds []int, t *testing.T, dbc *DBConnections) {
	for _, userId := range createdUserIds {
		err := dbc.UserService.DeleteUserAndSession(userId)
		if err != nil {
			t.Errorf("didn't expect error, got %v\n", err)
		}
	}

}
