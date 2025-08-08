package controllers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func HandlerExecuteTemplate(template ExecutorTemplate, fileName string, data any) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "text/html")
		template.ExecTemplate(w, r, fileName, data)
	}
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
