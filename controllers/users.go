package controllers

import (
	"net/http"
)

type Users struct {
	Templates struct {
		New ExcecutorTemplate
	}
}

func (u Users) ExecTemplate(w http.ResponseWriter, baseTemplate string, data any) {
	u.Templates.New.ExecTemplate(w, "signup.gohtml", nil)
}
