package models

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"time"
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

func (ss *SessionService) ExpirePreviousSessionsAndCreateNewSessionByUserId(userID int) (session *Session, err error) {
	err = ss.ExpireSessionsTokensByUserId(userID)
	if err != nil {
		return nil, err
	}
	session, err = ss.Create(userID)
	if err != nil {
		return nil, err
	}
	return session, nil
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
	returning id, user_id, token_hash;
	`, userID, hashedToken)
	returnedSession := &Session{}
	returnedSession.Token = newToken
	err = row.Scan(&returnedSession.ID, &returnedSession.UserID, &returnedSession.TokenHash)
	if err != nil {
		return nil, err
	}
	return returnedSession, nil
}
func (ss *SessionService) ExpireSessionsTokensByUserId(userID int) error {
	_, err := ss.db.Exec(`
		UPDATE sessions	
		SET is_expired =($1)
		WHERE user_id =($2);
	`, true, userID)
	if err != nil {
		return HandlePgError(err)
	}
	return nil
}

func (ss *SessionService) GetNonExpiredSessionsByUserId(userID int) (numSessions int, err error) {
	var count int
	row := ss.db.QueryRow(`
		SELECT COUNT(*)
		FROM sessions
		WHERE user_id=($1)
		AND is_expired=($2);
	`, userID, false)
	err = row.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (ss *SessionService) DeleteAllSessionsTokensByUserId(userID int) (err error) {
	_, err = ss.db.Exec(`
	DELETE from sessions
	WHERE user_id = ($1);
	`, userID)
	if err != nil {
		return err
	}
	return nil
}

func (ss *SessionService) CheckRequireRedirect(token string, cutOffTime time.Time) (isRequireRedirect bool) {
	type hashExpiryStruct struct {
		id        int
		expiresOn time.Time
	}
	var hashExpiry hashExpiryStruct
	tokenHash := HashSessionToken(token)
	row := ss.db.QueryRow(`
		SELECT id, expires_on
		FROM sessions
		WHERE token_hash=($1);
	`, tokenHash)
	// if any error occurs, then we take it that either no row was returned or there was an error, so we need to redirect
	if err := row.Scan(&hashExpiry.id, &hashExpiry.expiresOn); err != nil {
		fmt.Println("err occured: ", err)
		return true

	}
	if hashExpiry.expiresOn.UTC().Before(cutOffTime.UTC()) {
		return true
	}
	return false
}

// ViaToken retrieves the session that is tied to the session token sent with the request
func (ss *SessionService) ViaToken(token string) (*Session, error) {
	hash := HashSessionToken(token)
	row := ss.db.QueryRow(`
	SELECT id, user_id, token_hash 
	FROM sessions
	WHERE token_hash =($1);
	`, hash)
	var session Session
	// if any error occurs, then we take it that either no row was returned or there was an error, so we need to redirect
	err := row.Scan(&session.ID, &session.UserID, &session.TokenHash)
	if err != nil {
		return nil, err
	}
	return &session, nil
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
