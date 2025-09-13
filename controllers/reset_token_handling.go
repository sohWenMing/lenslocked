package controllers

import (
	"net/http"
)

func getTokenFromRequest(r *http.Request) (token string) {
	queryParams := r.URL.Query()
	token = queryParams.Get("token")
	return token
}
