package controllers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/sohWenMing/lenslocked/models"
)

func HandleSignupForm(dbc *models.DBConnections) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		emailAddress, password, err := parseEmailAndPasswordFromForm(r)
		if err != nil {
			http.Error(w, fmt.Sprintf("form could not be parsed: %s", err.Error()), http.StatusBadRequest)
			return
		}
		newUserToCreate := models.UserEmailToPlainTextPassword{
			Email:             emailAddress,
			PlainTextPassword: password,
		}
		user, err := dbc.UserService.CreateUser(newUserToCreate)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		sessionInformation := user.Session
		sessionToken := sessionInformation.Token
		SetSessionCookietoResponseWriter(sessionToken, w)
		http.Redirect(w, r, "/user/about", http.StatusFound)
	}
}

// closure function to allow access to the models.DBConnections type that returns a handler that can be used in main
// program

func HandleSignInForm(dbc *models.DBConnections) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		emailAddress, password, err := parseEmailAndPasswordFromForm(r)
		if err != nil {
			http.Error(w, fmt.Sprintf("form could not be parsed: %s", err.Error()), http.StatusBadRequest)
			return
		}
		userToPassword := models.UserEmailToPlainTextPassword{
			Email:             emailAddress,
			PlainTextPassword: password}

		loggedInUserInfo, err := dbc.UserService.LoginUser(userToPassword)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		sessionToken := loggedInUserInfo.Session.Token
		SetSessionCookietoResponseWriter(sessionToken, w)
		http.Redirect(w, r, "/user/about", http.StatusFound)
	}
}
func HandleForgotPasswordForm(dbc *models.DBConnections) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		email, err := ParseEmailFromForgetPasswordForm(r)
		if err != nil {
			http.Error(w, fmt.Sprintf("form could not be parsed: %s", err.Error()), http.StatusBadRequest)
			return
		}
		userInfo, err := dbc.UserService.GetUserByEmail(strings.ToLower(email))
		if err != nil {
			// TODO: Implement logging function
			fmt.Println("error: ", err)
		}
		fmt.Println("userInfo: ", userInfo)
		http.Redirect(w, r, "/check_email", http.StatusFound)
	}
}

func parseEmailAndPasswordFromForm(r *http.Request) (email, password string, err error) {
	err = r.ParseForm()
	if err != nil {
		return email, password, err
	}
	email = r.PostForm.Get("email")
	password = r.PostForm.Get("password")
	return email, password, nil
}

func ParseEmailFromForgetPasswordForm(r *http.Request) (email string, err error) {
	err = r.ParseForm()
	if err != nil {
		return email, err
	}
	email = r.PostForm.Get("email")
	fmt.Println("email entered: ", email)
	return email, nil
}
