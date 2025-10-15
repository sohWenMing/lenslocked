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
func (mwr *CookieAuthMWResult) SetIsErrOnExpireSessionByToken(input bool) {
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
			"167.172.78.219:80",
			"http://167.172.78.219",
			"http://167.172.78.219:80",
		},
		),
		csrf.Path("/"),
	)
}

func GetCSRFTokenFromRequest(r *http.Request) template.HTML {
	return csrf.TemplateField(r)
}

/*
CookieAuthMiddleWare returns a middleware that checks for the existence of a session cookie.

If a session cookie is not found, then it will redirect user to the sign in page.

if a session cookie is found but it is expired, will redirect user to to login page after expiring the session,
setting isExpired to true in database.

If a session cookie is found but not expired, will refresh the session adding 15 minutes to the session and set the UserId in the context
before passing on to the next middleware.

Writer passed in is used to record results from the middleware to be used for testing purposes, can set to nil for actual usage
isTestExpiry bool is used for testing purposes, can set to nil for actual usage
*/
func CookieAuthMiddleWare(ss *models.SessionService, writer io.Writer, isRedirect bool, isTestExpiry bool) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestTime := time.Now().UTC()

			// ##### For Testing ##### //
			// sets requestTime to 1 hour later than actual request time, to check forced expiry of session
			if isTestExpiry {
				requestTime = requestTime.Add(60 * time.Minute)
				fmt.Println("request time checked: ", requestTime)
			}

			//cookieAuthMWTResult used to record what happened in the middleware, used for testing purposes to write to writer
			cookieAuthMWRResult := &CookieAuthMWResult{}
			if writer != nil {
				defer helpers.WriteToWriter(writer, cookieAuthMWRResult)
			}

			//Attempts to get the sessionToken as a cookie, from the request
			sessionToken, isFound := GetSessionCookieFromRequest(r)

			// ##### For Testing ##### //
			cookieAuthMWRResult.SetIsCookieFoundFromGetSessionCookie(isFound)

			// if session cookie is not found, sets UserId in Context to 0, and moves on to next handler
			if !isFound {
				cookieAuthMWRResult.SetUserIdFromSession(0)
				r = setUserIdInContextForRequestZero(r)
				GoToPageOrRedirectToSignIn(isRedirect, next, w, r)
				return
			}
			// looks for session in the database, and checks expiry of session
			isSessionExpired, isSessionFound := ss.CheckSessionExpired(sessionToken, requestTime)

			// ##### For Testing ##### //
			cookieAuthMWRResult.SetIsSessionFoundInDatabase(isSessionFound)
			cookieAuthMWRResult.SetIsRedirectFromCheckSessionExpired(isSessionExpired)

			// if no session found, set the UserId to 0, and move on
			if !isSessionFound {
				cookieAuthMWRResult.SetUserIdFromSession(0)
				r = setUserIdInContextForRequestZero(r)
				GoToPageOrRedirectToSignIn(isRedirect, next, w, r)
				return
			}

			// if session time is expired, expire the session in the database, ser the UserId to 0, and move on
			if isSessionExpired {
				err := ss.ExpireSessionByToken(sessionToken)
				if err != nil {
					cookieAuthMWRResult.SetIsErrOnExpireSessionByToken(true)
					//TODO: implement logging function for error
					fmt.Println("error occured: ", err)
				} else {
					cookieAuthMWRResult.SetIsTokenSetToExpired(true)
				}
				cookieAuthMWRResult.SetUserIdFromSession(0)
				r = setUserIdInContextForRequestZero(r)
				GoToPageOrRedirectToSignIn(isRedirect, next, w, r)
				return
			}

			// else, attempt to refresh the session
			session, refreshErr := ss.RefreshSession(sessionToken, requestTime)

			// if error on refresh, set UserId to 0 in context, and move on
			if refreshErr != nil {
				cookieAuthMWRResult.SetIsErrorOnRefreshSession(true)
				//TODO: implement logging function for error
				cookieAuthMWRResult.SetUserIdFromSession(0)
				r = setUserIdInContextForRequestZero(r)
				GoToPageOrRedirectToSignIn(isRedirect, next, w, r)
				return
			}

			// session refreshed - setUserId in context to the returned userId from the session, and move on
			cookieAuthMWRResult.SetIsTokenSetToRefreshed(true)
			cookieAuthMWRResult.SetUserIdFromSession(session.UserID)
			ctx := context.WithValue(r.Context(), userIdKey, session.UserID)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)

		})
	}
}

/*
UserContext object houses methods that allow user related operations on request context
*/
type UserContext struct {
	userService *models.UserService
}

/*
NewUserContext returns pointer to a UserContext struct - which is a structure that has*
methods for getting and setting of user information on a request context
*/
func NewUserContext(us *models.UserService) *UserContext {
	return &UserContext{
		us,
	}
}

/*
SetUserMW takes the userId from the request context and checks whether or not the user
exists in the database. If the user does exist, then will move on the the next handler,
else will redirect to sign in.
*/
func (uc *UserContext) SetUserMW() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userId := r.Context().Value(userIdKey).(int)

			userInfo, _ := uc.userService.GetUserById(userId)
			ctx := context.WithValue(r.Context(), userInfoKey, userInfo)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

func (uc *UserContext) GetUserIdFromCtx(ctx context.Context) (userId int, isFound bool) {
	userId, ok := ctx.Value(getUserIdKey()).(int)
	if !ok {
		return 0, false
	}
	return userId, true
}
func (uc *UserContext) GetUserInfoFromCtx(ctx context.Context) (userInfo models.UserInfo, isFound bool) {
	userInfo, ok := ctx.Value(getUserInfoKey()).(models.UserInfo)
	if !ok {
		return models.UserInfo{}, false
	}
	return userInfo, true
}

func GetUserIdFromRequestContext(r *http.Request) (userId int, isFound bool) {
	ctx := r.Context()
	userId, ok := ctx.Value(getUserIdKey()).(int)
	if !ok {
		return 0, false
	}
	return userId, true
}

func GetUserInfoFromContext(r *http.Request) (userInfo models.UserInfo, isFound bool) {
	ctx := r.Context()
	userInfo, ok := ctx.Value(getUserInfoKey()).(models.UserInfo)
	if !ok {
		return models.UserInfo{}, false
	}
	return userInfo, true
}

func setUserIdInContextForRequestZero(r *http.Request) *http.Request {
	ctx := context.WithValue(r.Context(), userIdKey, 0)
	r = r.WithContext(ctx)
	return r
}

// func redirectToSignIn(w http.ResponseWriter, r *http.Request) {
// 	http.Redirect(w, r, "/signin", http.StatusFound)
// }

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
func getUserIdKey() contextKey {
	return userIdKey
}
func getUserInfoKey() contextKey {
	return userInfoKey
}
