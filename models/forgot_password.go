package models

import (
	"database/sql"
	"time"

	uuid "github.com/google/uuid"
)

type ForgotPWService struct {
	db *sql.DB
}

func (fpws *ForgotPWService) NewToken() (newToken uuid.UUID, err error) {
	newUUID := uuid.New()
	expires_on := time.Now().Add(time.Duration(15 * time.Minute)).UTC()
	row := fpws.db.QueryRow(
		`
		INSERT INTO forgot_password_tokens(token, expires_on)
		VALUES($1, $2)
		returning token;
		`, newUUID, expires_on,
	)
	var returnedToken uuid.UUID
	err = row.Scan(&returnedToken)
	if err != nil {
		return uuid.UUID{}, err
	}
	return returnedToken, nil
}
