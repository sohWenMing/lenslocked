package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/sohWenMing/lenslocked/models"
	"github.com/sohWenMing/lenslocked/services"
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
func HandleForgotPasswordForm(dbc *models.DBConnections, baseUrl string, emailer services.Emailer) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("baseUrl: ", baseUrl)
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
		newToken, err := dbc.ForgotPWService.NewToken(userInfo.ID)
		if err != nil {
			// TODO: Implement logging function
			fmt.Println("error: ", err)
		}
		fmt.Println("newToken returned: ", newToken)
		fmt.Println("userInfo: ", userInfo)
		urlToReturn := fmt.Sprintf("%s/reset_password?token=%s", baseUrl, newToken.String())
		fmt.Println("urlToReturn: ", urlToReturn)
		http.Redirect(w, r, "/check_email", http.StatusFound)
	}
}

func HandlerResetPasswordForm(dbc *models.DBConnections) func(w http.ResponseWriter, r *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("OK, the reset password form got submitted")

		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
		}
		fmt.Println("Form: ", r.Form)

		err = validatePasswordReset(r.Form)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("passwords must match"))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Reached reset handler form"))
	})
}

func validatePasswordReset(form url.Values) error {
	confirmPassword := form.Get("confirm-password")
	enterPassword := form.Get("enter-password")
	if enterPassword != confirmPassword {
		return errors.New("passwords must match")
	}
	return nil
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
