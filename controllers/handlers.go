package controllers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sohWenMing/lenslocked/models"
	"github.com/sohWenMing/lenslocked/views"
)

func GetHandlerFromTplMap(tplMap views.TplMap, fileName string) func(http.ResponseWriter, *http.Request) {
	template := tplMap[fileName]
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "text/html")
		template.Execute(w, nil)
	}
}
func GetHandlerforAbout(tplMap views.TplMap, fileName string) func(writer http.ResponseWriter, request *http.Request) {
	template := tplMap[fileName]
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("template: ", template)
		urlParam := chi.URLParam(r, "persona")
		fmt.Println("persona in handler: ", urlParam)
		if urlParam == "" {
			ErrNotFoundHandler(w, r)
			return
		}
		user, ok := models.UserMap[urlParam]
		if !ok {
			ErrNotFoundHandler(w, r)
			return
		}
		w.Header().Set("content-type", "text/html")
		template.Execute(w, user)
	}
}

func ErrNotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "404 not found")
}
