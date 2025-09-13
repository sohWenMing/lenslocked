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
func HandleForgotPasswordForm(dbc *models.DBConnections, baseUrl string, emailer *services.EmailService) func(w http.ResponseWriter, r *http.Request) {
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
			http.Redirect(w, r, "/check_email", http.StatusFound)
			return
		}
		newToken, err := dbc.ForgotPWService.NewToken(userInfo.ID)
		if err != nil {
			// TODO: Implement logging function
			fmt.Println("error: ", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		urlToReturn := fmt.Sprintf("%s/reset_password?token=%s", baseUrl, newToken.String())
		fmt.Println("urlToReturn: ", urlToReturn)

		emailData := services.EmailData{
			URL: urlToReturn,
		}

		emailBuf := bytes.Buffer{}
		err = emailer.EmailTemplate.EmailHTMLTpl.ExecuteTemplate(
			&emailBuf, "reset_password_email.gohtml", emailData,
		)
		if err != nil {
			fmt.Println("error: ", err)
			// TODO: Implement logging function
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
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
			fmt.Println("error: ", err)
			// TODO: Implement logging function
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/check_email", http.StatusFound)
	}
}

func HandlerResetPasswordForm(dbc *models.DBConnections) func(w http.ResponseWriter, r *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
		}

		err = validatePasswordReset(r.Form)
		if err != nil {
			fmt.Println("err in validatePasswordReset: ", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("passwords must match"))
			return
		}

		token, err := getForgotPWToken(r, dbc)
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		isValid := token.CheckIsValid()
		if !isValid {
			http.Error(w, "Token is no longer valid. Please make another request to reset password", http.StatusBadRequest)
			return
		}

		confirmedPassword := r.Form.Get("confirm-password")

		newHash, err := models.GenerateBcryptHash(confirmedPassword)
		if err != nil {
			http.Error(w, "Internal Error", http.StatusBadRequest)
			return
		}
		fmt.Println("newHash: ", newHash)

		err = dbc.ForgotPWService.DeleteForgetPasswordToken(token.UserId)
		if err != nil {
			http.Error(w, "Internal Error", http.StatusBadRequest)
			return
		}

		err = dbc.UserService.UpdatePasswordHash(token.UserId, newHash)
		if err != nil {
			http.Error(w, "Internal Error", http.StatusBadRequest)
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
