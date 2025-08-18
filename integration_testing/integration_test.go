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
	"sync"
	"testing"

	"github.com/sohWenMing/lenslocked/controllers"
	"github.com/sohWenMing/lenslocked/helpers"
	"github.com/sohWenMing/lenslocked/models"
)

var isDev = false
var dbc *models.DBConnections

type safeCounter struct {
	counter int
	mu      sync.Mutex
}

func (s *safeCounter) getIncrementedCounter() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.counter++
	return s.counter
}

var emailCounter = safeCounter{}

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

// test struct used in TestCookieAuthMiddleWare, extracted from test function
// to allow for type to be accessed by evalCookieAuthMWResult
type cookieAuthMWTest struct {
	name                                  string
	isTestBlankCookieInRequest            bool
	isTestRedirectFromCheckSessionExpired bool
	userInfo                              models.UserEmailToPlainTextPassword
}

func (c *cookieAuthMWTest) createUniqueEmail() {
	emailCounter := emailCounter.getIncrementedCounter()
	c.userInfo.Email = fmt.Sprintf("%s%d", c.userInfo.Email, emailCounter)
}

func evalCookieAuthMWResult(
	loggedInUserId int,
	test cookieAuthMWTest,
	cookieAuthMWResultIn *controllers.CookieAuthMWResult,
	t *testing.T) {
	//init working expected result
	expected := mapExpectedCookieAuthMWResult(test, loggedInUserId)
	if !reflect.DeepEqual(cookieAuthMWResultIn, expected) {
		t.Errorf("got %v\n, want %v\n", helpers.PrettyJSON(cookieAuthMWResultIn), helpers.PrettyJSON(expected))
		return
	}
}

func mapExpectedCookieAuthMWResult(test cookieAuthMWTest, userIdIn int) (expected *controllers.CookieAuthMWResult) {
	expected = &controllers.CookieAuthMWResult{}
	if test.isTestBlankCookieInRequest {
		expected.SetIsCookieFoundFromGetSessionCookie(false)
		expected.SetIsRedirectFromCheckSessionExpired(false)
		expected.SetIsSessionFoundInDatabase(false)
		expected.SetIsTokenSetToExpired(false)
		expected.SetIssErrOnExpireSessionByToken(false)
		expected.SetIsErrorOnRefreshSession(false)
		expected.SetIsTokenSetToRefreshed(false)
		expected.SetUserIdFromSession(0)

	} else if test.isTestRedirectFromCheckSessionExpired {
		expected.SetIsCookieFoundFromGetSessionCookie(true)
		expected.SetIsRedirectFromCheckSessionExpired(true)
		expected.SetIsSessionFoundInDatabase(true)
		expected.SetIsTokenSetToExpired(true)
		expected.SetIssErrOnExpireSessionByToken(false)
		expected.SetIsErrorOnRefreshSession(false)
		expected.SetIsTokenSetToRefreshed(false)
		expected.SetUserIdFromSession(0)
	} else {
		expected.SetIsCookieFoundFromGetSessionCookie(true)
		expected.SetIsRedirectFromCheckSessionExpired(false)
		expected.SetIsSessionFoundInDatabase(true)
		expected.SetIsTokenSetToExpired(false)
		expected.SetIssErrOnExpireSessionByToken(false)
		expected.SetIsErrorOnRefreshSession(false)
		expected.SetIsTokenSetToRefreshed(true)
		expected.SetUserIdFromSession(userIdIn)
	}
	return expected
}

func TestCookieAuthMiddleWare(t *testing.T) {

	tests := []cookieAuthMWTest{
		{
			"happy flow",
			false,
			false,
			models.UserEmailToPlainTextPassword{Email: "hello@test.com", PlainTextPassword: "Holoq123holoq123"},
		},
		{
			"test redirect from no session found",
			true,
			false,
			models.UserEmailToPlainTextPassword{Email: "hello@test.com", PlainTextPassword: "Holoq123holoq123"},
		},
		{
			"test redirect from expired session",
			false,
			true,
			models.UserEmailToPlainTextPassword{Email: "hello@test.com", PlainTextPassword: "Holoq123holoq123"},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.createUniqueEmail()
			//init
			createdUserIds := []int{}
			// cleanup action wrapped in closure, so that eval of createdUserIds will be delayed until actual end of overall function call
			defer func() {
				models.CleanUpCreatedUserIds(createdUserIds, t, dbc)
			}()

			//setup - create and login the user, so that we can get the session token to add as a cookie into the request

			createUserShouldReturn := createUser(t, test.userInfo, &createdUserIds)
			if createUserShouldReturn {
				return
			}
			loggedInUser, loginUserShouldReturn := loginUser(t, test.userInfo)
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
			mw := controllers.CookieAuthMiddleWare(dbc.SessionService, buf, true, test.isTestRedirectFromCheckSessionExpired)
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
			evalCookieAuthMWResult(loggedInUserId, test, cookieAuthMWResult, t)
		})
	}
}

type processSignOutTest struct {
	name                     string
	isTestNoSessionOnSignout bool
	expected                 controllers.ProcessSignoutResult
}

func evalProcessSIgnOutTest(test processSignOutTest,
	result controllers.ProcessSignoutResult,
	t *testing.T) {
	if test.isTestNoSessionOnSignout {
		test.expected.IsRedirectBecauseNoSession = true
	} else {
		test.expected.IsRedirectAfterExpiringSessionToken = true
		test.expected.IsSetExpireSessionCookie = true
	}
	if !reflect.DeepEqual(test.expected, result) {
		t.Errorf("got %s\n want %s\n", helpers.PrettyJSON(result),
			helpers.PrettyJSON(test.expected))
	}
}

func TestProcessSignOut(t *testing.T) {
	tests := []processSignOutTest{
		{
			"happy flow, redirect should be happening after expiring of session token",
			false,
			controllers.ProcessSignoutResult{},
		},
		{
			"test no session attached to request",
			true,
			controllers.ProcessSignoutResult{},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			createdUserIds := []int{}

			//cleanup function
			defer func() {
				models.CleanUpCreatedUserIds(createdUserIds, t, dbc)
			}()
			userInfo := models.UserEmailToPlainTextPassword{Email: fmt.Sprintf("hello@test.com%d", emailCounter.getIncrementedCounter()),
				PlainTextPassword: "Holoq123holoq123"}
			shouldReturn := createUser(t, userInfo, &createdUserIds)
			if shouldReturn {
				return
			}
			//log in the user, so that we can get the session token tied to the user
			loggedInUser, shouldReturn := loginUser(t, userInfo)
			if shouldReturn {
				return
			}
			req, err := http.NewRequest(http.MethodGet, "/", nil)
			if !test.isTestNoSessionOnSignout {
				sessionToken := loggedInUser.Session.Token
				sessionCookie := controllers.MapSessionCookie(sessionToken)
				req.AddCookie(sessionCookie)
			}
			if err != nil {
				t.Errorf("didn't expect error, got %v", err)
				return
			}
			responseRecorder := httptest.NewRecorder()
			buf := &bytes.Buffer{}
			controllers.HandlerSignOut(dbc.SessionService, buf)(responseRecorder, req)

			var processSignOutResult controllers.ProcessSignoutResult
			_ = json.Unmarshal(buf.Bytes(), &processSignOutResult)
			evalProcessSIgnOutTest(test, processSignOutResult, t)
		})
	}
}

func AddCookieToRequest(loggedInUser *models.UserIdToSession, newRequest *http.Request) {
	sessionToken := loggedInUser.Session.Token
	sessionCookie := controllers.MapSessionCookie(sessionToken)
	newRequest.AddCookie(sessionCookie)
}

func loginUser(t *testing.T, userInfo models.UserEmailToPlainTextPassword) (*models.UserIdToSession, bool) {
	loggedInUser, err := dbc.UserService.LoginUser(userInfo)
	if err != nil {
		t.Errorf("didn't expect error, got %v\n", err)
		return nil, true
	}
	return loggedInUser, false
}

func createUser(t *testing.T, userInfo models.UserEmailToPlainTextPassword, createdUserIds *[]int) bool {
	createdUser, err := dbc.UserService.CreateUser(userInfo)
	*createdUserIds = append(*createdUserIds, createdUser.ID)
	if err != nil {
		t.Errorf("didn't expect error, got %v\n", err)
		return true
	}
	return false
}
