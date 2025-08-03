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
	initFormData := views.SignUpSignInFormData
	initFormData.SetEmailValue(r.FormValue("email"))
	s.FormController.Templates.ExecTemplate(w, "signup.gohtml", initFormData)
}

func InitSignupFormController(template ExecutorTemplate) *SignupFormController {
	return &SignupFormController{
		FormController: FormController{
			template,
		},
	}
}

// ##### SignIn Form Controller Definition #####
type SignInFormController struct {
	FormController FormController
}

func (s *SignInFormController) Load(w http.ResponseWriter, r *http.Request) {
	initFormData := views.SignUpSignInFormData
	initFormData.SetEmailValue(r.FormValue("email"))
	s.FormController.Templates.ExecTemplate(w, "signin.gohtml", initFormData)
}

func InitSignInFormController(template ExecutorTemplate) *SignInFormController {
	return &SignInFormController{
		FormController: FormController{
			template,
		},
	}
}

func HandleSignupForm(dbc *models.DBConnections) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, fmt.Sprintf("form could not be parsed: %s", err.Error()), http.StatusBadRequest)
			return
		}
		emailAddress := r.PostForm.Get("email")
		password := r.PostForm.Get("password")
		user, err := dbc.UserService.CreateUser(emailAddress, password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "<p>user with email %s has been successfully created", user.Email)
		fmt.Fprintf(w, "<p>email address: %s</p>", emailAddress)
		fmt.Fprintf(w, "<p>password: %s</p>", password)
	}
}
