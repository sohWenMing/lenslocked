package controllers

import (
	"fmt"
	"net/http"

	"github.com/sohWenMing/lenslocked/models"
	"github.com/sohWenMing/lenslocked/views"
)

/*
What am i trying to achieve
I want to unify all the different forms, under a type , so that each of them can load. The Loader should be a http.Handler
the thing to take into account is that each loader, will have a different set of things that it needs to do

so ....
if we take the Load method, that will be the unifier, and then from there we can deal with each thing that we need
*/

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
		fmt.Println("handler managed to get to end of parsing of form ")
		user, err := dbc.UserService.CreateUser(emailAddress, password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Println("handler managed to get to end of creating the user")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "<p>user with email %s has been successfully created", user.Email)
		fmt.Fprintf(w, "<p>email address: %s</p>", emailAddress)
		fmt.Fprintf(w, "<p>password: %s</p>", password)
	}
}

// ##### Practice Form Controller Definition #####
type PraticeFormController struct {
	FormController FormController
}

func InitPracticeFormController(template ExecutorTemplate) *PraticeFormController {
	return &PraticeFormController{
		FormController: FormController{
			template,
		},
	}
}

func (s *PraticeFormController) Load(w http.ResponseWriter, r *http.Request) {
	initFormData := views.PracticeFormData
	s.FormController.Templates.ExecTemplate(w, "practice_form.gohtml", initFormData)
}
func HandlePracticeForm(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, fmt.Sprintf("form could not be parsed: %s", err.Error()), http.StatusBadRequest)
		return
	}
	fmt.Println("full form ", r.PostForm)
	firstName := r.PostForm.Get("first_name")
	lastName := r.PostForm.Get("last_name")
	isChecked := r.PostForm.Get("testCheckBox")
	fmt.Printf("isChecked: \"%s\n\"", isChecked)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "<p>first_name: %s</p>", firstName)
	fmt.Fprintf(w, "<p>last_name: %s</p>", lastName)
	if isChecked == "isChecked" {
		fmt.Fprint(w, "form was checked upon submit")
	}
}
