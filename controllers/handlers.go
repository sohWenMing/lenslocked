package controllers

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sohWenMing/lenslocked/helpers"
	"github.com/sohWenMing/lenslocked/models"
	"github.com/sohWenMing/lenslocked/views"
)

func HandlerExecuteTemplate(template ExecutorTemplateWithCSRF, fileName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId, isFound := GetUserIdFromRequestContext(r)
		if !isFound {
			fmt.Println("userId not found")
		} else {
			fmt.Println("userId: ", userId)
		}
		otherPageData, err := views.BaseTemplatesToData.GetDataForTemplate(fileName)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		pageData := views.PageData{
			UserId:    userId,
			OtherData: otherPageData,
		}
		if fileName == "signup.gohtml" || fileName == "sign.gohtml" {
			signInSignUpFormData, ok := otherPageData.(views.SignInSignUpForm)
			if !ok {
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			}
			updatedSignInSignUpFormData := setSignInSignUpFormData(r, signInSignUpFormData)
			pageData.OtherData = updatedSignInSignUpFormData
		}

		fmt.Println("pageData: ", pageData)

		csrfToken := GetCSRFTokenFromRequest(r)
		w.Header().Set("content-type", "text/html")
		template.ExecTemplateWithCSRF(w, r, csrfToken, fileName, pageData)
	}
}

func setSignInSignUpFormData(r *http.Request, signInSignUpFormData views.SignInSignUpForm) views.SignInSignUpForm {
	signInSignUpFormData.SetEmailValue(r.FormValue("email"))
	return signInSignUpFormData
}

func GetUrlParam(r *http.Request, param string) (returnedString string, err error) {
	returnedString = chi.URLParam(r, param)
	if returnedString == "" {
		return "", errors.New("param could not be found")
	}
	return returnedString, nil
}

func ErrNotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "404 not found")
}

func TestHandler(testText string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, testText)
	}
}

func TestSendCookie(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("email")
	if err != nil {
		fmt.Println("err: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	returnedString := fmt.Sprintf("name: %s value: %s", cookie.Name, cookie.Value)
	fmt.Fprint(w, returnedString)
	fmt.Fprintf(w, "Headers %v\n", r.Header)
}

type ProcessSignoutResult struct {
	IsRedirectBecauseNoSession          bool
	IsErrOnExpireSessionToken           bool
	IsRedirectAfterExpiringSessionToken bool
	IsSetExpireSessionCookie            bool
}

func (p *ProcessSignoutResult) SetIsRedirectBecauseNoSession(bool) {
	p.IsRedirectBecauseNoSession = true
}
func (p *ProcessSignoutResult) SetIsErrOnExpireSessionToken(bool) {
	p.IsErrOnExpireSessionToken = true
}
func (p *ProcessSignoutResult) SetIsRedirectAfterExpiringSessionToken(bool) {
	p.IsRedirectAfterExpiringSessionToken = true
}
func (p *ProcessSignoutResult) SetIsSetExpireSessionCookie(bool) {
	p.IsSetExpireSessionCookie = true
}

// Processes a sign out request - writer that is passed in should be used for testing purposes. Set nil to writer for actual application
func HandlerSignOut(ss *models.SessionService, writer io.Writer) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var result ProcessSignoutResult
		if writer != nil {
			defer func() {
				helpers.WriteToWriter(writer, result)
			}()
		}
		token, isFound := GetSessionCookieFromRequest(r)
		if !isFound {
			result.SetIsRedirectBecauseNoSession(true)
			http.Redirect(w, r, "/signin", http.StatusFound)
			return
		}
		SetExpireSessionCookieToResponseWriter(token, w)
		result.SetIsSetExpireSessionCookie(true)
		err := ss.ExpireSessionByToken(token)
		if err != nil {
			result.SetIsErrOnExpireSessionToken(true)
			// TODO: implement logging function
			fmt.Println(err)
		}
		result.SetIsRedirectAfterExpiringSessionToken(true)
		http.Redirect(w, r, "/signin", http.StatusFound)
	})
}
