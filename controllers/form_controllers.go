package controllers

import (
	"fmt"
	"net/http"

	"github.com/sohWenMing/lenslocked/models"
	"github.com/sohWenMing/lenslocked/views"
)

type FormLoader interface {
	Load(w http.ResponseWriter, r *http.Request)
}

type FormNameToLoader map[string]FormLoader

// map so that main program can retrieve the correct form loader depending on handler

func InitFormNameToLoader(template ExecutorTemplateWithCSRF) FormNameToLoader {
	return FormNameToLoader{
		"signup_form": InitSignupFormController(template),
		"signin_form": InitSignInFormController(template),
	}
}

//used to return the FormNameToLoader map that will be used

type FormController struct {
	Templates ExecutorTemplateWithCSRF
}

// ##### Signup Form Controller Definition #####
/*

The Signin and Signup form controllers are defined as different types - this is so that while the Load method
can make both the Signup and SignIn form controllers fufil the FormLoader interface, there is individual control
over each controller

in this way there is individual control over the data that is passed in to the ExecTemplate function, and also
which template will be eventually be used as the base template

*/
type SignupFormController struct {
	FormController FormController
}

func (s *SignupFormController) Load(w http.ResponseWriter, r *http.Request) {
	initFormData := setSignInSignUpFormData(r)
	csrfToken := GetCSRFTokenFromRequest(r)
	s.FormController.Templates.ExecTemplateWithCSRF(w, r, csrfToken, "signup.gohtml", initFormData)
}

func InitSignupFormController(template ExecutorTemplateWithCSRF) *SignupFormController {
	return &SignupFormController{
		FormController: FormController{
			template,
		},
	}
}

func HandleSignupForm(dbc *models.DBConnections) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		emailAddress, password, err := parseEmailAndPasswordFromForm(r)
		if err != nil {
			http.Error(w, fmt.Sprintf("form could not be parsed: %s", err.Error()), http.StatusBadRequest)
			return
		}
		newUserToCreate := models.UserToPlainTextPassword{
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

// ##### SignIn Form Controller Definition #####
type SignInFormController struct {
	FormController FormController
}

func (s *SignInFormController) Load(w http.ResponseWriter, r *http.Request) {
	initFormData := setSignInSignUpFormData(r)
	csrfToken := GetCSRFTokenFromRequest(r)
	s.FormController.Templates.ExecTemplateWithCSRF(w, r, csrfToken, "signin.gohtml", initFormData)
}

func InitSignInFormController(template ExecutorTemplateWithCSRF) *SignInFormController {
	return &SignInFormController{
		FormController: FormController{
			template,
		},
	}
}

func HandleSignInForm(dbc *models.DBConnections) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		emailAddress, password, err := parseEmailAndPasswordFromForm(r)
		if err != nil {
			http.Error(w, fmt.Sprintf("form could not be parsed: %s", err.Error()), http.StatusBadRequest)
			return
		}
		userToPassword := models.UserToPlainTextPassword{
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

func parseEmailAndPasswordFromForm(r *http.Request) (email, password string, err error) {
	err = r.ParseForm()
	if err != nil {
		return email, password, err
	}
	email = r.PostForm.Get("email")
	password = r.PostForm.Get("password")
	return email, password, nil
}

func setSignInSignUpFormData(r *http.Request) views.SignInSignUpForm {
	initFormData := views.SignUpSignInFormData
	initFormData.SetEmailValue(r.FormValue("email"))
	return initFormData
}
