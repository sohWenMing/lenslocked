package controllers

import "net/http"

type ExecutorTemplate interface {
	ExecTemplate(w http.ResponseWriter, baseTemplate string, data any)
	// so in this way, anything that has the ExecTemplate function filfills this interface
}
