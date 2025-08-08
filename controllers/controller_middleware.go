package controllers

import (
	"html/template"
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

func GetCSRFTokenFromRequest(r *http.Request) template.HTML {
	return csrf.TemplateField(r)
}

func CookieAuthMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("email")
		if err != nil {
			http.Redirect(w, r, "/signin", http.StatusPermanentRedirect)
			return
		}
		if cookie.Value == "" {
			http.Redirect(w, r, "/login", http.StatusPermanentRedirect)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// the middleware takes in the next handler - so if everything passes then it will delegate on to the next handler
// if .Use is set from the router, then it will when delegate to the router's handler
