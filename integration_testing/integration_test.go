package integrationtesting

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/sohWenMing/lenslocked/controllers"
	"github.com/sohWenMing/lenslocked/models"
)

var isDev = false
var dbc *models.DBConnections

type cookieAuthMWTest struct {
	name                                  string
	isTestBlankCookieInRequest            bool
	isTEstRedirectFromCheckSessionExpired bool
	expiry                                time.Time
	userInfo                              models.UserToPlainTextPassword
}

func TestMain(m *testing.M) {
	envVars, err := models.LoadEnv("../.env")
	if err != nil {
		log.Fatal(err)
	}
	isDev = envVars.IsDev
	databaseConnection, err := models.InitDBConnections()
	if err != nil {
		fmt.Println("error occured during initialisation of db connection during test")
		os.Exit(1)
	}
	dbc = databaseConnection
	code := m.Run()
	os.Exit(code)
}

func TestEnvLoading(t *testing.T) {
	want := true
	got := isDev
	if got != want {
		t.Errorf("Got %v, want %v\n", got, want)
	}
}

func TestCookieAuthMiddleWare(t *testing.T) {

	tests := []cookieAuthMWTest{
		{
			"happy flow",
			false,
			false,
			time.Now(),
			models.UserToPlainTextPassword{Email: "hello@test.com", PlainTextPassword: "Holoq123holoq123"},
		},
		{
			"test redirect from no session found",
			true,
			false,
			time.Now(),
			models.UserToPlainTextPassword{Email: "hello@test.com", PlainTextPassword: "Holoq123holoq123"},
		},
		{
			"test redirect from expired session",
			false,
			true,
			time.Now().Add(60 * time.Minute),
			models.UserToPlainTextPassword{Email: "hello@test.com", PlainTextPassword: "Holoq123holoq123"},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			//init
			createdUserIds := []int{}
			// cleanup action wrapped in closure, so that eval of createdUserIds will be delayed until actual end of overall function call
			defer func() {
				models.CleanUpCreatedUserIds(createdUserIds, t, dbc)
			}()

			//setup - create and login the user, so that we can get the session token to add as a cookie into the request
			createUserShouldReturn := createUser(t, test, &createdUserIds)
			if createUserShouldReturn {
				return
			}
			loggedInUser, loginUserShouldReturn := loginUser(t, test)
			if loginUserShouldReturn {
				return
			}
			loggedInUserId := loggedInUser.UserID
			//setup - just a test handler, so that we can wrap it with the CookieAuthMiddleware to see how the middleware responds
			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("OK"))
			})
			/*
				setup - defining the instance of the cookieAuthMiddleware, passing in a buffer so that we can read from the buffer later
				for evals. testHandler gets wrapped in the middleware, which will result in the eventual handler that will handle the
				reqeuest setup in the test
			*/
			buf := &bytes.Buffer{}
			mw := controllers.CookieAuthMiddleWare(dbc.SessionService, buf, test.expiry)
			wrappedHandler := mw(testHandler)

			newRequest, err := http.NewRequest(http.MethodGet, "/test", nil)
			if err != nil {
				t.Errorf("didn't expect error, got %v\n", err)
				return
			}
			if !test.isTestBlankCookieInRequest {
				AddCookieToRequest(loggedInUser, newRequest)
			}
			requestRecorder := httptest.NewRecorder()
			wrappedHandler.ServeHTTP(requestRecorder, newRequest)
			cookieAuthMWResult := &controllers.CookieAuthMWResult{}
			unMarshalErr := json.Unmarshal(buf.Bytes(), cookieAuthMWResult)
			if unMarshalErr != nil {
				t.Errorf("didn't expect error, got %v\n", err)
				return
			}
			fmt.Println("cookieAuthMWResult: ", cookieAuthMWResult.PrettyJSON())
			evalCookieAuthMWResult(test, loggedInUserId, cookieAuthMWResult, t)
		})
	}
}

func AddCookieToRequest(loggedInUser *models.User, newRequest *http.Request) {
	sessionToken := loggedInUser.Session.Token
	sessionCookie := controllers.MapSessionCookie(sessionToken)
	newRequest.AddCookie(sessionCookie)
}

func loginUser(t *testing.T, test cookieAuthMWTest) (*models.User, bool) {
	loggedInUser, err := dbc.UserService.LoginUser(test.userInfo)
	if err != nil {
		t.Errorf("didn't expect error, got %v\n", err)
		return nil, true
	}
	return loggedInUser, false
}

func createUser(t *testing.T, test cookieAuthMWTest, createdUserIds *[]int) bool {
	createdUser, err := dbc.UserService.CreateUser(test.userInfo)
	*createdUserIds = append(*createdUserIds, createdUser.ID)
	if err != nil {
		t.Errorf("didn't expect error, got %v\n", err)
		return true
	}
	return false
}

func evalCookieAuthMWResult(
	test cookieAuthMWTest,
	loggedInUserId int,
	cookieAuthMWRequestIn *controllers.CookieAuthMWResult,
	t *testing.T) {
	//init working expected result
	expected := &controllers.CookieAuthMWResult{}

	if test.isTestBlankCookieInRequest {
		expected.SetIsRedirectFromGetSessionCookie(true)
		//TODO: Handle this case later
	} else if test.isTEstRedirectFromCheckSessionExpired {
		expected.SetIsTokenSetToExpired(true)
		expected.SetIsRedirectFromCheckSessionExpired(true)
		expected.SetIsSessionFound(true)

	} else {
		expected.SetIsSessionFound(true)
		expected.SetIsTokenSetToRefreshed(true)
		expected.SetUserIdFromSession(loggedInUserId)
	}
	if !reflect.DeepEqual(cookieAuthMWRequestIn, expected) {
		t.Errorf("got %v, want %v\n", cookieAuthMWRequestIn, expected)
		return
	}

}
