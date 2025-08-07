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

func InitFormNameToLoader(template ExecutorTemplate) FormNameToLoader {
	return FormNameToLoader{
		"signup_form": InitSignupFormController(template),
		"signin_form": InitSignInFormController(template),
	}
}

type FormController struct {
	Templates ExecutorTemplate
}

// ##### Signup Form Controller Definition #####
type SignupFormController struct {
	FormController FormController
}

func (s *SignupFormController) Load(w http.ResponseWriter, r *http.Request) {
	initFormData := setSignInSignUpFormData(r)
	s.FormController.Templates.ExecTemplate(w, "signup.gohtml", initFormData)
}

func InitSignupFormController(template ExecutorTemplate) *SignupFormController {
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
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "<p>user with email %s has been successfully created", user.Email)
	}
}

// ##### SignIn Form Controller Definition #####
type SignInFormController struct {
	FormController FormController
}

func (s *SignInFormController) Load(w http.ResponseWriter, r *http.Request) {
	initFormData := setSignInSignUpFormData(r)
	s.FormController.Templates.ExecTemplate(w, "signin.gohtml", initFormData)
}

func InitSignInFormController(template ExecutorTemplate) *SignInFormController {
	return &SignInFormController{
		FormController: FormController{
			template,
		},
	}
}

func HandlerSigninForm(dbc *models.DBConnections) func(w http.ResponseWriter, r *http.Request) {
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
		cookie := http.Cookie{
			Name:     "email",
			Value:    loggedInUserInfo.Email,
			Path:     "/",
			HttpOnly: true,
		}
		http.SetCookie(w, &cookie)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "<p>user with email %s has been successfully logged in", loggedInUserInfo.Email)
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
	// initFormData.SetEmailValue(r.FormValue("email"))
	return initFormData
}
