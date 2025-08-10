package models

import (
	"fmt"
	"os"
	"testing"
)

var dbc *DBConnections
var userService UserService

func TestMain(m *testing.M) {
	databaseConnection, err := InitDBConnections()
	if err != nil {
		fmt.Println("error occured during initialisation of db connection during test")
		os.Exit(1)
	}
	dbc = databaseConnection
	code := m.Run()
	if err := dbc.DB.Close(); err != nil {
		fmt.Println("error closing db: ", err)
	}
	os.Exit(code)
}

func TestCreateUser(t *testing.T) {
	createdUserIds := []int{}
	type test struct {
		name           string
		enteredInfo    UserToPlainTextPassword
		expectedErrMsg string
		isErrExpected  bool
	}

	tests := []test{
		{
			"test email validation",
			UserToPlainTextPassword{"hello@test.com", "Holoq123holoq123"},
			"",
			false,
		},
		{
			"test duplicate email ",
			UserToPlainTextPassword{"hello@test.com", "Holoq123holoq123"},
			"",
			true,
		},
		{
			"test wrong password ",
			UserToPlainTextPassword{"hello1@test.com", "12345"},
			"",
			true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			user, err := dbc.UserService.CreateUser(test.enteredInfo)
			if err == nil {
				createdUserIds = append(createdUserIds, user.ID)
			}
			switch test.isErrExpected {
			case true:
				if err == nil {
					t.Error("error expected, didn't get one\n")
					return
				}
			default:
				if err != nil {
					t.Errorf("didn't expect error, got %v\n", err)
					return
				}
				if user.ID == 0 {
					t.Errorf("user was not created")
					return
				}
				if user.Session.ID == 0 {
					t.Errorf("session was not created")
				}
			}
		})
	}
	// cleanup
	cleanupCreatedUserIds(createdUserIds, t)
}

func TestLoginUser(t *testing.T) {
	createdUserIds := []int{}
	type test struct {
		name                 string
		isErrExpectedOnLogin bool
		userInfo             UserToPlainTextPassword
	}
	tests := []test{
		{
			"test happy flow",
			false,
			UserToPlainTextPassword{
				"hello@test.com",
				"Holoq123holoq123",
			},
		},
		{
			"test failed login",
			true,
			UserToPlainTextPassword{
				"hello1@test.com",
				"Holoq123holoq123",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			createdUser, err := dbc.UserService.CreateUser(test.userInfo)
			if err != nil {
				t.Errorf("didn't expect error, got %v\n", err)
				return
			}
			createdUserIds = append(createdUserIds, createdUser.ID)

			switch test.isErrExpectedOnLogin {
			case true:
				changedUserInfo := UserToPlainTextPassword{
					test.userInfo.Email, "fail_password",
				}
				_, err := dbc.UserService.LoginUser(changedUserInfo)
				if err == nil {
					t.Errorf("expected error, didn't get one")
				}
			default:
				loggedInUser, err := dbc.UserService.LoginUser(test.userInfo)
				if err != nil {
					t.Errorf("didn't expect error, got %v\n", err)
				}
				if loggedInUser.Session.UserID != createdUser.ID {
					t.Errorf("got session userId %d want session userId %d", loggedInUser.Session.UserID, createdUser.ID)
				}
			}
		})
	}
	// cleanup
	cleanupCreatedUserIds(createdUserIds, t)
}

func cleanupCreatedUserIds(createdUserIds []int, t *testing.T) {
	for _, userId := range createdUserIds {
		err := dbc.UserService.DeleteUserAndSession(userId)
		if err != nil {
			t.Errorf("didn't expect error, got %v\n", err)
		}
	}

}
