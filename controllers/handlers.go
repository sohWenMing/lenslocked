package controllers

import (
	"fmt"
	"net/http"

	"github.com/sohWenMing/lenslocked/views"
)

func GetHandlerFromTplMap(tplMap views.TplMap, fileName string) func(http.ResponseWriter, *http.Request) {
	template := tplMap[fileName]
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "text/html")
		template.Execute(w, nil)
	}
}

func ErrNotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "404 not found")
}
