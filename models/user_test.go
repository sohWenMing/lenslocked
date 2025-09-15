package models

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/sohWenMing/lenslocked/helpers"
)

func TestCreateUser(t *testing.T) {
	createdUserIds := []int{}
	type test struct {
		name           string
		enteredInfo    UserEmailToPlainTextPassword
		expectedErrMsg string
		isErrExpected  bool
	}

	tests := []test{
		{
			"test email validation",
			baseUserEmailToPlainTextPassword,
			"",
			false,
		},
		{
			"test duplicate email ",
			baseUserEmailToPlainTextPassword,
			"",
			true,
		},
		{
			"test wrong password ",
			UserEmailToPlainTextPassword{"hello1@test.com", "12345"},
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
					return
				}
			}
		})
	}
	// cleanup
	CleanUpCreatedUserIds(createdUserIds, t, dbc)
}

func TestLoginUser(t *testing.T) {
	createdUserIds := []int{}
	type test struct {
		name                 string
		isErrExpectedOnLogin bool
		userInfo             UserEmailToPlainTextPassword
	}
	tests := []test{
		{
			"test happy flow",
			false,
			UserEmailToPlainTextPassword{
				"hello@test.com",
				"Holoq123holoq123",
			},
		},
		{
			"test failed login",
			true,
			UserEmailToPlainTextPassword{
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
				changedUserInfo := UserEmailToPlainTextPassword{
					test.userInfo.Email, "fail_password",
				}
				_, err := dbc.UserService.LoginUser(changedUserInfo)
				if err == nil {
					t.Errorf("expected error, didn't get one")
					return
				}
			default:
				loggedInUser, err := dbc.UserService.LoginUser(test.userInfo)
				if err != nil {
					t.Errorf("didn't expect error, got %v\n", err)
					return
				}
				if loggedInUser.Session.UserID != createdUser.ID {
					t.Errorf("got session userId %d want session userId %d", loggedInUser.Session.UserID, createdUser.ID)
					return
				}
			}
		})
	}
	// cleanup
	CleanUpCreatedUserIds(createdUserIds, t, dbc)
}

func TestExpireSessionsByUserId(t *testing.T) {
	createdUserIds := []int{}
	type test struct {
		name     string
		userInfo UserEmailToPlainTextPassword
	}
	tests := []test{
		{
			"test happy flow expire sessions by user id",
			UserEmailToPlainTextPassword{
				"hello@test.com",
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
			loggedInUser, err := dbc.UserService.LoginUser(test.userInfo)
			if err != nil {
				t.Errorf("didn't expect error, got %v\n", err)
				return
			}
			err = dbc.UserService.ExpireSessionsTokensByUserId(loggedInUser.ID)
			if err != nil {
				t.Errorf("didn't expect error, got %v\n", err)
				return
			}
			nonExpiredCount, err := dbc.UserService.GetNonExpiredSessionsByUserId(loggedInUser.ID)
			if err != nil {
				t.Errorf("didn't expect error, got %v\n", err)
				return
			}
			if nonExpiredCount != 0 {
				t.Errorf("expected nonExpiredCount %d, got %d", 0, nonExpiredCount)
			}
		})
	}
	// cleanup
	CleanUpCreatedUserIds(createdUserIds, t, dbc)
}

func TestDeleteUser(t *testing.T) {
	type test struct {
		name     string
		userInfo UserEmailToPlainTextPassword
	}
	tests := []test{
		{
			"test delete user",
			UserEmailToPlainTextPassword{
				"hello@test.com",
				"Holoq123holoq123",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := dbc.UserService.CreateUser(test.userInfo)
			if err != nil {
				t.Errorf("didn't expect error, got %v\n", err)
				return
			}
			loggedInUser, err := dbc.UserService.LoginUser(test.userInfo)
			if err != nil {
				t.Errorf("didn't expect error, got %v\n", err)
				return
			}
			nonExpiredCount, err := dbc.UserService.GetNonExpiredSessionsByUserId(loggedInUser.ID)
			if err != nil {
				t.Errorf("didn't expect error, got %v\n", err)
				return
			}
			if nonExpiredCount != 1 {
				t.Errorf("expected nonExpiredCount %d, got %d", 1, nonExpiredCount)
			}
			err = dbc.UserService.DeleteUserAndSession(loggedInUser.ID)
			if err != nil {
				t.Errorf("didn't expect error, got %v\n", err)
				return
			}
			userCount, err := dbc.UserService.GetUserCountById(loggedInUser.ID)
			if err != nil {
				t.Errorf("didn't expect error, got %v\n", err)
				return
			}
			if userCount != 0 {
				t.Errorf("expected userCount %d, got %d", 0, userCount)
			}
			nonExpiredCountAfterDelete, err := dbc.UserService.GetNonExpiredSessionsByUserId(loggedInUser.ID)
			if err != nil {
				t.Errorf("didn't expect error, got %v\n", err)
				return
			}
			if nonExpiredCountAfterDelete != 0 {
				t.Errorf("expected nonExpiredCountAfterDelete %d, got %d", 0, nonExpiredCountAfterDelete)
			}
		})
	}

}

func TestLogoutUser(t *testing.T) {
	createdUserIds := []int{}
	type test struct {
		name     string
		userInfo UserEmailToPlainTextPassword
	}
	tests := []test{
		{
			"test happy flow logout",
			UserEmailToPlainTextPassword{
				"hello@test.com",
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
			loggedInUser, err := dbc.UserService.LoginUser(test.userInfo)
			if err != nil {
				t.Errorf("didn't expect error, got %v\n", err)
				return
			}
			err = dbc.UserService.LogoutUser(loggedInUser.ID)
			if err != nil {
				t.Errorf("didn't expect error, got %v\n", err)
				return
			}
			nonExpiredSessionCount, err := dbc.UserService.GetNonExpiredSessionsByUserId(loggedInUser.ID)
			if err != nil {
				t.Errorf("didn't expect error, got %v\n", err)
				return
			}
			if nonExpiredSessionCount != 0 {
				t.Errorf("expected nonExpiredSessionCount %d, got %d", 0, nonExpiredSessionCount)
			}
		})
	}

	// cleanup
	CleanUpCreatedUserIds(createdUserIds, t, dbc)
}

func TestRequireRedirect(t *testing.T) {
	createdUserIds := []int{}
	type test struct {
		name             string
		isExpectRedirect bool
		userInfo         UserEmailToPlainTextPassword
	}
	tests := []test{
		{
			"happy flow, should not require redirect",
			false,
			UserEmailToPlainTextPassword{
				"hello@test.com",
				"Holoq123holoq123",
			},
		},
		{
			"should require redirect",
			true,
			UserEmailToPlainTextPassword{
				"hello@test1.com",
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
			loggedInUser, err := dbc.UserService.LoginUser(test.userInfo)
			if err != nil {
				t.Errorf("didn't expect error, got %v\n", err)
				return
			}
			token := loggedInUser.Token
			var isRequireRedirect bool

			switch test.isExpectRedirect {
			case true:
				fmt.Println("test case ran with isExpectRedirect to true")
				isRequireRedirect, _ = dbc.SessionService.CheckSessionExpired(token, time.Now().Add(16*time.Minute))
			default:
				isRequireRedirect, _ = dbc.SessionService.CheckSessionExpired(token, time.Now())
			}

			if isRequireRedirect != test.isExpectRedirect {
				t.Errorf("got %t, want %t\n", isRequireRedirect, test.isExpectRedirect)
			}
		})
	}
	// cleanup
	CleanUpCreatedUserIds(createdUserIds, t, dbc)

}

func TestGetUserById(t *testing.T) {
	type test struct {
		name              string
		wantString        string
		expectedErrMsg    string
		CreatedUserInputs UserEmailToPlainTextPassword
		isErrExpected     bool
		want              UserInfo
	}

	createdUserIds := []int{}
	tests := []test{
		{
			"happy flow, able to find created user",
			"",
			"",
			baseUserEmailToPlainTextPassword,
			false,
			UserInfo{0, strings.ToLower(baseUserEmailToPlainTextPassword.Email)},
		},
		{
			"testing userId that does not exist",
			"",
			"No user could be found with that user id",
			baseUserEmailToPlainTextPassword,
			true,
			UserInfo{0, strings.ToLower(baseUserEmailToPlainTextPassword.Email)},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			defer func() {
				CleanUpCreatedUserIds(createdUserIds, t, dbc)
			}()
			createdUser, err := dbc.UserService.CreateUser(test.CreatedUserInputs)
			if err != nil {
				t.Errorf("didn't expect error, got %v", err)
				return
			}
			createdUserIds = append(createdUserIds, createdUser.ID)
			switch test.isErrExpected {
			case false:
				returnedUser, err := dbc.UserService.GetUserById(createdUser.ID)
				if err != nil {
					t.Errorf("didn't expect error, got %v", err)
					return
				}
				test.want.ID = returnedUser.ID
				if !reflect.DeepEqual(test.want, returnedUser) {
					t.Errorf("got %s\n want %s\n", helpers.PrettyJSON(returnedUser), helpers.PrettyJSON(test.want))
					return
				}
				test.wantString = fmt.Sprintf("UserId: %d Email: %s", test.want.ID, test.CreatedUserInputs.Email)
				if test.wantString != returnedUser.String() {
					t.Errorf("got %s\n want %s\n",
						fmt.Sprintf(`"%s"`, returnedUser.String()),
						fmt.Sprintf(`"%s"`, test.wantString),
					)
				}
			case true:
				_, err := dbc.UserService.GetUserById(createdUser.ID + 1)
				if err == nil {
					t.Errorf("expected error, didn't get one")
					return
				}
				if err.Error() != test.expectedErrMsg {
					t.Errorf("got errMsg %s\n want errMsg %s\n", err.Error(), test.expectedErrMsg)
				}
				var handledError *HandledError
				if !errors.As(err, &handledError) {
					t.Errorf("error that was returned was not a handledError type")
				}
				if !errors.Is(handledError.err, sql.ErrNoRows) {
					t.Errorf("err in handledError was not of type sql.ErrNoRows: %v", handledError.err)
				}
			}
		})
	}
}
