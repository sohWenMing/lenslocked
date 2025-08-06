package views

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/sohWenMing/lenslocked/models"
)

type Template struct {
	htmlTpl *template.Template
}

type TplMap map[string]Template

var tplStrings = []string{
	"home.gohtml",
	"contact.gohtml",
	"faq.gohtml",
	"persona.gohtml",
	"tailwind_widgets.gohtml",
	"signup.gohtml",
	"signin.gohtml",
	"practice_form.gohtml",
	"test_cookie.gohtml",
}

var BaseTemplateToData = map[string]any{
	"home.gohtml":          nil,
	"contact.gohtml":       nil,
	"faq.gohtml":           models.QuestionsToAnswers,
	"signup.gohtml":        SignUpSignInFormData,
	"signin.gohtml":        SignUpSignInFormData,
	"practice_form.gohtml": nil,
	"test_cookie.gohtml":   nil,
}

//go:embed templates/*
var FS embed.FS

func (t *Template) ExecTemplate(w http.ResponseWriter, baseTemplate string, data any) {
	err := t.htmlTpl.ExecuteTemplate(w, baseTemplate, data)
	if err != nil {
		log.Printf("parsing template: %v", err)
		http.Error(w,
			"There was an error parsing the template",
			http.StatusInternalServerError)
		return
	}
}

func LoadTemplates() (tpl *Template) {
	tpl = &Template{}

	templateStrings := getTemplatePaths(tplStrings, "templates")
	LoadedTemplate := TemplateMust(template.ParseFS(FS, templateStrings...))
	tpl.htmlTpl = LoadedTemplate
	return tpl
}

func TemplateMust(t *template.Template, err error) *template.Template {
	if err != nil {
		panic(err)
	}
	return t
}

func getTemplatePaths(tplStrings []string, baseFolderName string) []string {
	fullPaths := make([]string, len(tplStrings))
	for i, tplString := range tplStrings {
		fullPath := fmt.Sprintf("%s/%s", baseFolderName, tplString)
		fullPaths[i] = fullPath
	}
	return fullPaths
}
