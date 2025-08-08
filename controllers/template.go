package controllers

import (
	"html/template"
	"net/http"
)

type ExecutorTemplate interface {
	ExecTemplate(w http.ResponseWriter, r *http.Request, baseTemplate string, data any)
	// so in this way, anything that has the ExecTemplate function filfills this interface
}
type ExecutorTemplateWithCSRF interface {
	ExecTemplateWithCSRF(w http.ResponseWriter, r *http.Request, csrfToken template.HTML, baseTemplate string, data any)
}
