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
		// cleanup
		fmt.Println("rcreated userIds:", createdUserIds)
	}
	for _, userId := range createdUserIds {
		err := dbc.UserService.DeleteUserAndSession(userId)
		if err != nil {
			t.Errorf("didn't expect error, got %v\n", err)
		}
	}
}
