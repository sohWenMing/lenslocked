package controllers

import (
	"net/http"

	"github.com/google/uuid"
)

func getTokenFromRequest(r *http.Request) (token string) {
	queryParams := r.URL.Query()
	token = queryParams.Get("token")
	return token
}

func parseTokenStringToUUID(token string) (parsedUUID uuid.UUID, err error) {
	parsedUUID, err = uuid.Parse(token)
	if err != nil {
		return uuid.UUID{}, err
	}
	return parsedUUID, nil
}
