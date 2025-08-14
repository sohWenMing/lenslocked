package controllers

import (
	"net/http"
)

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
