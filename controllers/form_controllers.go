package controllers

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/uuid"
	"github.com/sohWenMing/lenslocked/models"
	"github.com/sohWenMing/lenslocked/services"
)

func HandleSignupForm(
	dbc *models.DBConnections,
	render func(w http.ResponseWriter, r *http.Request, fileName string, errorMsgs []string),
) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		emailAddress, password, err := parseEmailAndPasswordFromForm(r)
		if err != nil {
			render(w, r, "signup.gohtml", []string{"form could not be parsed. please reload, and try again"})
			return
		}
		newUserToCreate := models.UserEmailToPlainTextPassword{
			Email:             emailAddress,
			PlainTextPassword: password,
		}
		user, err := dbc.UserService.CreateUser(newUserToCreate)
		if err != nil {
			render(w, r, "signup.gohtml", []string{err.Error()})
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

func HandleSignInForm(dbc *models.DBConnections,
	render func(w http.ResponseWriter, r *http.Request, fileName string, errorMsgs []string),
) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		emailAddress, password, err := parseEmailAndPasswordFromForm(r)
		if err != nil {
			render(w, r, "signin.gohtml", []string{"form could not be parsed. please reload, and try again"})
			return
		}
		userToPassword := models.UserEmailToPlainTextPassword{
			Email:             emailAddress,
			PlainTextPassword: password}

		loggedInUserInfo, err := dbc.UserService.LoginUser(userToPassword)

		if err != nil {
			render(w, r, "signin.gohtml", []string{"there was a problem with the username and password. please check and try again"})
			return
		}
		sessionToken := loggedInUserInfo.Session.Token
		SetSessionCookietoResponseWriter(sessionToken, w)
		http.Redirect(w, r, "/galleries/list", http.StatusFound)
	}
}
func HandleForgotPasswordForm(dbc *models.DBConnections, baseUrl string, emailer *services.EmailService,
	render func(w http.ResponseWriter, r *http.Request, fileName string, errorMsgs []string),
) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		email, err := ParseEmailFromForgetPasswordForm(r)
		if err != nil {
			render(w, r, "forgot_password.gohtml", []string{"form could not be parsed. please reload, and try again"})
			return
		}
		userInfo, err := dbc.UserService.GetUserByEmail(strings.ToLower(email))
		if err != nil {
			render(w, r, "forgot_password.gohtml", []string{"No user exists with that email. Please try again"})
			return
		}
		newToken, err := dbc.ForgotPWService.NewToken(userInfo.ID)
		if err != nil {
			render(w, r, "forgot_password.gohtml", []string{"There was a problem with the request. Please try again."})
			return
		}
		urlToReturn := fmt.Sprintf("%s/reset_password?token=%s", baseUrl, newToken.String())

		emailData := services.EmailData{
			URL: urlToReturn,
		}

		emailBuf := bytes.Buffer{}
		err = emailer.EmailTemplate.EmailHTMLTpl.ExecuteTemplate(
			&emailBuf, "reset_password_email.gohtml", emailData,
		)
		if err != nil {
			render(w, r, "forgot_password.gohtml", []string{"There was a problem with the request. Please try again."})
			return
		}

		err = emailer.SendMail(services.Email{
			From:        "wenming.soh@gmail.com",
			To:          email,
			Content:     emailBuf.String(),
			ContentType: "text/html",
			Cc:          []string{},
		}, nil)

		if err != nil {
			render(w, r, "forgot_password.gohtml", []string{"There was a problem sending the email. Please try again in a while."})
			return
		}
		http.Redirect(w, r, "/check_email", http.StatusFound)
	}
}

func HandlerResetPasswordForm(dbc *models.DBConnections,
	render func(w http.ResponseWriter, r *http.Request, fileName string, errorMsgs []string),
) func(w http.ResponseWriter, r *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		err := r.ParseForm()
		if err != nil {
			render(w, r, "reset_password.gohtml", []string{"there was an error parsing the form"})
			return
		}

		err = validatePasswordReset(r.Form)
		if err != nil {
			render(w, r, "reset_password.gohtml", []string{"passwords must match"})
			return
		}

		token, err := getForgotPWToken(r, dbc)
		if err != nil {
			render(w, r, "reset_password.gohtml", []string{"this link has expired - please make a new request."})
			return
		}

		isValid := token.CheckIsValid()
		if !isValid {
			render(w, r, "reset_password.gohtml", []string{"this link has expired - please make a new request."})
			return
		}

		confirmedPassword := r.Form.Get("confirm-password")

		newHash, err := models.GenerateBcryptHash(confirmedPassword)
		if err != nil {
			render(w, r, "reset_password.gohtml", []string{"there was an internal error - please try again and contact support if the problem persists."})
			return
		}

		err = dbc.ForgotPWService.DeleteForgetPasswordToken(token.UserId)
		if err != nil {
			render(w, r, "reset_password.gohtml", []string{"there was an internal error - please try again and contact support if the problem persists."})
			return
		}

		err = dbc.UserService.UpdatePasswordHash(token.UserId, newHash)
		if err != nil {
			render(w, r, "reset_password.gohtml", []string{"there was an internal error - please try again and contact support if the problem persists."})
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Password has been reset, please login"))
	})
}

func getForgotPWToken(r *http.Request, dbc *models.DBConnections) (models.ForgotPasswordToken, error) {
	var forgotPasswordToken models.ForgotPasswordToken
	token := r.Form.Get("forgot_password_token")
	tokenUUID, err := uuid.Parse(token)
	if err != nil {
		return forgotPasswordToken, nil
	}
	forgotPasswordToken, err = dbc.ForgotPWService.GetForgotPWToken(tokenUUID)
	return forgotPasswordToken, err
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
