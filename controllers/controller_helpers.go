package controllers

import (
	"net/http"
)

// checks if the cookie with key "sessionToken" can be found. if not found or val is blank, isMustRedirect will return true
func GetSessionCookieFromRequest(r *http.Request) (token string, isFound bool) {
	sessionCookie, err := r.Cookie("sessionToken")
	if err != nil {
		return "", false
	}
	if sessionCookie.Value == "" {
		return "", false
	}
	token = sessionCookie.Value
	return token, true
}
