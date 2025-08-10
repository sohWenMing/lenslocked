package controllers

import (
	"context"
	"html/template"
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/sohWenMing/lenslocked/models"
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

func CookieAuthMiddleWare(ss *models.SessionService) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sessionCookie, err := r.Cookie("sessionToken")
			if err != nil {
				http.Redirect(w, r, "/signin", http.StatusFound)
				return
			}
			if sessionCookie.Value == "" {
				http.Redirect(w, r, "/signin", http.StatusFound)
				return
			}

			session, err := ss.ViaToken(sessionCookie.Value)
			if err != nil {
				http.Redirect(w, r, "/signin", http.StatusFound)
			}
			ctx := context.WithValue(r.Context(), "userId", session.UserID)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

// the middleware takes in the next handler - so if everything passes then it will delegate on to the next handler
// if .Use is set from the router, then it will when delegate to the router's handler
