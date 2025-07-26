package controllers

import (
	"fmt"
	"net/http"

	"github.com/sohWenMing/lenslocked/views"
)

func HandlerExecuteTemplate(template views.Template, fileName string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "text/html")
		template.ExecTemplate(w, fileName)
	}
}

func ErrNotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "404 not found")
}
