package models

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
)

const SessionTokenSize = 32

type Session struct {
	ID     int
	UserID int
	/*
		Token is only set when creating a new session. When looking up a session
		this will be left empty, as we only store the hash of a session token
		in our database, and we can't reverse it into a raw token
	*/
	Token     string
	TokenHash string
}

type SessionService struct {
	db *sql.DB
}

/*
Create will create a new session for the user provided. The session token will
be created as the Token field on the Session type, but only the hashed session
token will be stored on the database.
*/
func (ss *SessionService) Create(userID int) (*Session, error) {
	newToken, err := SessionToken()
	if err != nil {
		return nil, err
	}
	hashedToken := HashSessionToken(newToken)
	row := ss.db.QueryRow(`
	INSERT into sessions(user_id, token_hash)
	VALUES($1, $2)
	returning id, user_id, token_hash 
	`, userID, hashedToken)
	returnedSession := &Session{}
	returnedSession.Token = newToken
	err = row.Scan(&returnedSession.ID, &returnedSession.UserID, &returnedSession.TokenHash)
	if err != nil {
		return nil, err
	}
	return returnedSession, nil
}

func (ss *SessionService) ViaToken(token string) (*Session, error) {
	//TODO, create the querying from the database, to get the session of the user
	return nil, nil
}

func VerifySessionToken(token string, hash string) (isVerified bool, err error) {
	if token == "" {
		return false, errors.New("token passed in cannot be nil")
	}
	if hash == "" {
		return false, errors.New("hash passed in cannot be nil")
	}
	hashedSessionToken := HashSessionToken(token)
	if hashedSessionToken != hash {
		return false, nil
	}
	return true, nil
}

func HashSessionToken(token string) string {
	h := sha256.New()
	h.Write([]byte(token))
	return base64.URLEncoding.EncodeToString(h.Sum(nil))
}

func SessionToken() (string, error) {
	return bytestring(SessionTokenSize)
}

func bytestring(n int) (string, error) {
	bytes, err := randBytes(n)
	if err != nil {
		return "", fmt.Errorf("string: %w", err)
	}
	returnedString := base64.URLEncoding.EncodeToString(bytes)
	return returnedString, nil
}

func randBytes(numBytes int) ([]byte, error) {
	bytes := make([]byte, numBytes)
	nRead, nErr := rand.Read(bytes)
	if nErr != nil {
		return nil, fmt.Errorf("bytes: %w", nErr)
	}
	if nRead < numBytes {
		return nil, fmt.Errorf("bytes: didn't read enough bytes")
	}
	return bytes, nil
}
