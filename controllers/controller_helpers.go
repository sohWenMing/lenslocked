package controllers

import (
	"net/http"
)

// checks if the cookie with key "sessionToken" can be found. if not found or val is blank, isMustRedirect will return true
func GetSessionCookieFromRequest(r *http.Request) (token string, isMustRedirect bool) {
	sessionCookie, err := r.Cookie("sessionToken")
	if err != nil {
		return "", true
	}
	if sessionCookie.Value == "" {
		return "", true
	}
	token = sessionCookie.Value
	return token, false
}
