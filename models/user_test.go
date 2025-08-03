package models

import (
	"errors"
	"os"
	"reflect"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

var dbc *DBConnections
var userService UserService

func TestMain(m *testing.M) {
	dbc = InitDBConnections()
	userService = UserService{
		dbc.DB,
	}
	defer dbc.DB.Close()
	code := m.Run()
	os.Exit(code)
}

func TestLoginUser(t *testing.T) {
	type test struct {
		testName       string
		userToPassword UserToPlainTextPassword
		expected       LoggedInUserInfo
		expectedErr    error
		isErrExpected  bool
	}

	tests := []test{
		{
			"first test, to get valid user",
			UserToPlainTextPassword{
				"wenming.soh@gmail.com",
				"Holoq123holoq123",
			},
			LoggedInUserInfo{
				1, "wenming.soh@gmail.com",
			},
			nil,
			false,
		},
		{
			"second test, should fail password",
			UserToPlainTextPassword{
				"wenming.soh@gmail.com",
				"failing_pw",
			},
			LoggedInUserInfo{},
			bcrypt.ErrMismatchedHashAndPassword,
			true,
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			user, err := userService.LoginUser(test.userToPassword)
			switch test.isErrExpected {
			case true:
				if err == nil {
					t.Errorf("expected error, didn't get one")
				}
				var target *HandledError
				if errors.As(err, &target) {
					if target.err != test.expectedErr {
						t.Errorf("expected error %v, got %v", test.expectedErr, target.err)
					}
				}

			default:
				if err != nil {
					t.Fatalf("didn't expect error, got err: %v", err)
				}
			}
			if !reflect.DeepEqual(user, test.expected) {
				t.Fatalf("got %v, want %v", user, test.expected)
			}
		})
	}
}
