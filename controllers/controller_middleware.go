package controllers

import (
	"context"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"time"

	"github.com/gorilla/csrf"
	"github.com/sohWenMing/lenslocked/helpers"
	"github.com/sohWenMing/lenslocked/models"
)

type CookieAuthMWResult struct {
	IsRedirectFromGetSessionCookie    bool
	IsRedirectFromCheckSessionExpired bool
	IsSessionFound                    bool
	IsErrOnExpireSessionByToken       bool
	IsTokenSetToExpired               bool
	IsErrOnRefreshSession             bool
	IsTokenSetToRefreshed             bool
	UserIdFromSession                 int
}

func (mwr *CookieAuthMWResult) SetIsRedirectFromGetSessionCookie(input bool) {
	mwr.IsRedirectFromGetSessionCookie = input
}
func (mwr *CookieAuthMWResult) SetIsRedirectFromCheckSessionExpired(input bool) {
	mwr.IsRedirectFromCheckSessionExpired = input
}
func (mwr *CookieAuthMWResult) SetIsSessionFound(input bool) {
	mwr.IsSessionFound = input
}
func (mwr *CookieAuthMWResult) SetIsTokenSetToExpired(input bool) {
	mwr.IsTokenSetToExpired = input
}
func (mwr *CookieAuthMWResult) SetIssErrOnExpireSessionByToken(input bool) {
	mwr.IsErrOnExpireSessionByToken = input
}
func (mwr *CookieAuthMWResult) SetIsErrorOnRefreshSession(input bool) {
	mwr.IsErrOnRefreshSession = input
}
func (mwr *CookieAuthMWResult) SetIsTokenSetToRefreshed(input bool) {
	mwr.IsTokenSetToRefreshed = input
}
func (mwr *CookieAuthMWResult) SetUserIdFromSession(userId int) {
	mwr.UserIdFromSession = userId
}

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

/*
Function checks for the existence of a session cookie, if does not exist, will redirect to signin page
writer passed in is to allow flexibility for capturing values during testing, in actual running code
should be set to nil
*/
func CookieAuthMiddleWare(ss *models.SessionService, writer io.Writer, expiry time.Time) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		requestTime := time.Now().UTC()
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			//cookieAuthMWTResult used to record what happened in the middleware, used for testing purposes to write to writer
			cookieAuthMWRResult := &CookieAuthMWResult{}
			if writer != nil {
				defer helpers.WriteToWriter(writer, cookieAuthMWRResult)
			}
			//checks for existence of session token in the cookies from the requesst
			token, isRequireRedirect := GetSessionCookieFromRequest(r)

			//for testing - writes to the cookieAuthMWRResult
			cookieAuthMWRResult.SetIsRedirectFromGetSessionCookie(isRequireRedirect)

			if isRequireRedirect {
				http.Redirect(w, r, "/signin", http.StatusFound)
				return
			}
			isRequireRedirect, isSessionFound := ss.CheckSessionExpired(token, expiry)
			cookieAuthMWRResult.SetIsSessionFound(isSessionFound)

			//for testing - writes to the cookieAuthMWRResult
			cookieAuthMWRResult.SetIsRedirectFromCheckSessionExpired(isRequireRedirect)

			if isRequireRedirect {
				if isSessionFound {
					err := ss.ExpireSessionByToken(token)
					if err != nil {
						cookieAuthMWRResult.SetIssErrOnExpireSessionByToken(true)
						//TODO: implement logging function for error
						fmt.Println("error occured: ", err)
					} else {
						cookieAuthMWRResult.SetIsTokenSetToExpired(true)
					}
				}
				http.Redirect(w, r, "/signin", http.StatusFound)
				return
			}
			session, refreshErr := ss.RefreshSession(token, requestTime)
			if refreshErr != nil {
				cookieAuthMWRResult.SetIsErrorOnRefreshSession(true)
				//TODO: implement logging function for error
				fmt.Println("error occured: ", refreshErr)
				http.Redirect(w, r, "/signin", http.StatusFound)
			} else {
				cookieAuthMWRResult.SetIsTokenSetToRefreshed(true)
			}
			cookieAuthMWRResult.SetUserIdFromSession(session.UserID)
			ctx := context.WithValue(r.Context(), "userId", session.UserID)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

// the middleware takes in the next handler - so if everything passes then it will delegate on to the next handler
// if .Use is set from the router, then it will when delegate to the router's handler
