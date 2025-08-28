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
	IsCookieFoundFromGetSessionCookie bool
	IsRedirectFromCheckSessionExpired bool
	IsSessionFoundInDatabase          bool
	IsErrOnExpireSessionByToken       bool
	IsTokenSetToExpired               bool
	IsErrOnRefreshSession             bool
	IsTokenSetToRefreshed             bool
	UserIdFromSession                 int
}

type contextKey string

const userIdKey = contextKey("userId")
const userInfoKey = contextKey("userInfo")

func (mwr *CookieAuthMWResult) SetIsCookieFoundFromGetSessionCookie(input bool) {
	mwr.IsCookieFoundFromGetSessionCookie = input
}
func (mwr *CookieAuthMWResult) SetIsRedirectFromCheckSessionExpired(input bool) {
	mwr.IsRedirectFromCheckSessionExpired = input
}
func (mwr *CookieAuthMWResult) SetIsSessionFoundInDatabase(input bool) {
	mwr.IsSessionFoundInDatabase = input
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
CookieAuthMiddleWare returns a middleware that checks for the existence of a session cookie.
If a session cookie is not found, then it will redirect user to the sign in page
if a session cookie is found but it is expired, will redirect user to to login page after expiring the session, setting isExpired to true in database
If a session coookie is found but not expired, will refresh the session adding 15 minutes
Writer passed in is used to record results from the middleware to be used for testing purposes, can set to nil for actual usage
isTestExpiry bool is used for testing purposes, can set to nil for actual usage
*/

// the middleware takes in the next handler - so if everything passes then it will delegate on to the next handler
// if .Use is set from the router, then it will when delegate to the router's handler

func CookieAuthMiddleWare(ss *models.SessionService, writer io.Writer, isRedirect bool, isTestExpiry bool) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestTime := time.Now().UTC()
			//testing purposes

			if isTestExpiry {
				requestTime = requestTime.Add(60 * time.Minute)
				fmt.Println("request time checked: ", requestTime)
			}

			//cookieAuthMWTResult used to record what happened in the middleware, used for testing purposes to write to writer
			cookieAuthMWRResult := &CookieAuthMWResult{}
			if writer != nil {
				defer helpers.WriteToWriter(writer, cookieAuthMWRResult)
			}
			token, isFound := GetSessionCookieFromRequest(r)

			//for testing - writes to the cookieAuthMWRResult
			cookieAuthMWRResult.SetIsCookieFoundFromGetSessionCookie(isFound)

			if !isFound {
				cookieAuthMWRResult.SetUserIdFromSession(0)
				r = setUserIdInContextForRequestZero(r)
				GoToPageOrRedirectToSignIn(isRedirect, next, w, r)
				return
			}
			isSessionExpired, IsSessionFoundInDatabase := ss.CheckSessionExpired(token, requestTime)

			//for testing - writes to the cookieAuthMWRResult
			cookieAuthMWRResult.SetIsSessionFoundInDatabase(IsSessionFoundInDatabase)
			cookieAuthMWRResult.SetIsRedirectFromCheckSessionExpired(isSessionExpired)

			if isSessionExpired {
				if IsSessionFoundInDatabase {
					err := ss.ExpireSessionByToken(token)
					if err != nil {
						cookieAuthMWRResult.SetIssErrOnExpireSessionByToken(true)
						//TODO: implement logging function for error
						fmt.Println("error occured: ", err)
					} else {
						cookieAuthMWRResult.SetIsTokenSetToExpired(true)
					}
				}
				cookieAuthMWRResult.SetUserIdFromSession(0)
				r = setUserIdInContextForRequestZero(r)
				GoToPageOrRedirectToSignIn(isRedirect, next, w, r)
				return
			}
			session, refreshErr := ss.RefreshSession(token, requestTime)
			if refreshErr != nil {
				cookieAuthMWRResult.SetIsErrorOnRefreshSession(true)
				//TODO: implement logging function for error
				cookieAuthMWRResult.SetUserIdFromSession(0)
				r = setUserIdInContextForRequestZero(r)
				GoToPageOrRedirectToSignIn(isRedirect, next, w, r)
				return
			} else {
				cookieAuthMWRResult.SetIsTokenSetToRefreshed(true)
			}
			cookieAuthMWRResult.SetUserIdFromSession(session.UserID)
			ctx := context.WithValue(r.Context(), userIdKey, session.UserID)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

func UserInfoMiddleWare(us *models.UserService) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			userId := r.Context().Value(userIdKey).(int)
			userInfo, err := us.GetUserById(userId)
			if err != nil {
				//TODO: Implement logging function
				fmt.Println(err)
			}
			ctx := context.WithValue(r.Context(), userInfoKey, userInfo)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

func setUserIdInContextForRequestZero(r *http.Request) *http.Request {
	ctx := context.WithValue(r.Context(), userIdKey, 0)
	r = r.WithContext(ctx)
	return r
}

func GoToPageOrRedirectToSignIn(isRedirect bool, next http.Handler, w http.ResponseWriter, r *http.Request) {
	switch isRedirect {
	case false:
		next.ServeHTTP(w, r)
		return
	default:
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}
}
func GetUserIdKey() contextKey {
	return userIdKey
}
func GetUserInfoKey() contextKey {
	return userInfoKey
}

func GetUserIdFromRequestContext(r *http.Request) (userId int, isFound bool) {
	ctx := r.Context()
	userId, ok := ctx.Value(GetUserIdKey()).(int)
	if !ok {
		return 0, false
	}
	return userId, true
}

func GetUserInfoFromContext(r *http.Request) (userInfo models.UserIdToEmail, isFound bool) {
	ctx := r.Context()
	userInfo, ok := ctx.Value(GetUserInfoKey()).(models.UserIdToEmail)
	if !ok {
		return models.UserIdToEmail{}, false
	}
	return userInfo, true
}
