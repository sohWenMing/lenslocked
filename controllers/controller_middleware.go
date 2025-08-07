package controllers

import (
	"net/http"

	"github.com/gorilla/csrf"
)

const secretKey = "gFvi45R4fy5xNBlnEeZtQbfAVCYEIAUX"

func CSRFProtect(isDev bool) func(http.Handler) http.Handler {
	return csrf.Protect([]byte(secretKey), csrf.Secure(isDev))
}
