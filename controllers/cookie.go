package controllers

import (
	"net/http"
)

/*
maps a cookie based on the session token that is created when a user logs in, and attaches it to a
http.ResponseWriter so cookie will be send in the response
*/
func SetSessionCookietoResponseWriter(sessionToken string, w http.ResponseWriter) {
	cookie := MapSessionCookie(sessionToken)
	// fmt.Println("cookie in SetSessionCookieToResponseWriter: ", cookie)
	http.SetCookie(w, cookie)
}
func SetExpireSessionCookieToResponseWriter(sessionToken string, w http.ResponseWriter) {
	cookie := MapExpireSessionCookie(sessionToken)
	// fmt.Println("cookie in SetExpireSessionCookieToResponseWriter: ", cookie)
	http.SetCookie(w, cookie)
}

func MapSessionCookie(token string) *http.Cookie {
	return mapCookie("sessionToken", token, "/", true, 15)
}

func MapExpireSessionCookie(token string) *http.Cookie {
	return mapCookie("sessionToken", token, "/", true, -1)
}

func mapCookie(name, value, path string, HTTPOnly bool, maxAgeInMinutes int) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     path,
		HttpOnly: HTTPOnly,
		MaxAge:   maxAgeInMinutes * 60,
	}
}
