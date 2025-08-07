package controllers

import (
	"net/http"

	"github.com/gorilla/csrf"
)

func CSRFProtect(isDev bool, secretKey string) func(http.Handler) http.Handler {
	isSetSecure := !isDev
	return csrf.Protect([]byte(secretKey),
		csrf.Secure(isSetSecure),
		csrf.TrustedOrigins([]string{
			"localhost:3000",
		}))
}
