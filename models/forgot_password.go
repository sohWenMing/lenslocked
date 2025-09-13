package models

import (
	"database/sql"
	"fmt"
	"time"

	uuid "github.com/google/uuid"
)

type ForgotPWService struct {
	db *sql.DB
}

type ForgotPasswordToken struct {
	Id        int
	UserId    int
	Token     uuid.UUID
	ExpiresOn time.Time
}

func (fpwt ForgotPasswordToken) GetExpiry() time.Time {
	return fpwt.ExpiresOn
}

func (fpwt ForgotPasswordToken) CheckIsValid() bool {
	expiry := fpwt.GetExpiry()
	return expiry.Compare(time.Now()) == +1
}

func (fpws *ForgotPWService) DeleteForgetPasswordToken(userId int) error {
	_, err := fpws.db.Exec(
		`
		DELETE FROM forgot_password_tokens
		WHERE user_id = ($1);
		`, userId)
	if err != nil {
		return err
	}
	return nil

}

func (fpws *ForgotPWService) NewToken(userId int) (newToken uuid.UUID, err error) {
	newUUID := uuid.New()
	expires_on := time.Now().Add(time.Duration(15 * time.Minute)).UTC()
	row := fpws.db.QueryRow(
		`
		INSERT INTO forgot_password_tokens(user_id, token, expires_on)
		VALUES($1, $2, $3)
		returning token;
		`, NoUserWithUserId, newUUID, expires_on,
	)
	var returnedToken uuid.UUID
	err = row.Scan(&returnedToken)
	if err != nil {
		return uuid.UUID{}, err
	}
	return returnedToken, nil
}

func (fpws *ForgotPWService) GetForgotPWToken(token uuid.UUID) (ForgotPasswordToken, error) {
	fmt.Println("token: ", token)
	row := fpws.db.QueryRow(
		`
		SELECT id, user_id, token, expires_on FROM forgot_password_tokens
		WHERE token=($1)
		`, token,
	)
	var fpwToken ForgotPasswordToken
	err := row.Scan(&fpwToken.Id, &fpwToken.UserId, &fpwToken.Token, &fpwToken.ExpiresOn)
	if err != nil {
		return ForgotPasswordToken{}, err
	}
	return fpwToken, nil
}
