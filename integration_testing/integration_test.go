package integrationtesting

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/sohWenMing/lenslocked/controllers"
	"github.com/sohWenMing/lenslocked/models"
)

var isDev = false
var dbc *models.DBConnections

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

	type test struct {
		name               string
		userInfo           models.UserToPlainTextPassword
		expectedTestResult controllers.CookieAuthMWResult
	}
	tests := []test{
		{
			"happy flow",
			models.UserToPlainTextPassword{Email: "hello@test.com", PlainTextPassword: "Holoq123holoq123"},
			controllers.CookieAuthMWResult{},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			createdUserIds := []int{}
			defer func() {
				models.CleanUpCreatedUserIds(createdUserIds, t, dbc)
			}()
			createdUser, err := dbc.UserService.CreateUser(test.userInfo)
			createdUserIds = append(createdUserIds, createdUser.ID)
			fmt.Println("createdUserIds: ", createdUserIds)
			if err != nil {
				t.Errorf("didn't expect error, got %v\n", err)
				return
			}
			loggedInUser, err := dbc.UserService.LoginUser(test.userInfo)
			if err != nil {
				t.Errorf("didn't expect error, got %v\n", err)
				return
			}
			sessionToken := loggedInUser.Session.Token
			sessionCookie := controllers.MapSessionCookie(sessionToken)
			newRequest, err := http.NewRequest(http.MethodGet, "/test", nil)
			if err != nil {
				t.Errorf("didn't expect error, got %v\n", err)
				return
			}
			newRequest.AddCookie(sessionCookie)
			buf := &bytes.Buffer{}
			mw := controllers.CookieAuthMiddleWare(dbc.SessionService, buf)
			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("OK"))
			})
			wrappedHandler := mw(testHandler)
			requestRecorder := httptest.NewRecorder()
			wrappedHandler.ServeHTTP(requestRecorder, newRequest)
			cookieAuthMWResult := &controllers.CookieAuthMWResult{}
			unMarshalErr := json.Unmarshal(buf.Bytes(), cookieAuthMWResult)
			if unMarshalErr != nil {
				t.Errorf("didn't expect error, got %v\n", err)
				return
			}
			fmt.Println("cookieAuthMWResult: ", cookieAuthMWResult)
		})

	}

}
