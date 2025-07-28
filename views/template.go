package views

import (
	"embed"
	"errors"
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
	"persona_multiple.gohtml",
	"tailwind.gohtml",
	"tailwind_widgets.gohtml",
	"signup.gohtml",
}

var BaseTemplateToData = map[string]any{
	"home.gohtml":             nil,
	"contact.gohtml":          nil,
	"faq.gohtml":              models.QuestionsToAnswers,
	"persona_multiple.gohtml": models.GetAllUsers(),
	"signup.gohtml":           SignUpFormData,
}

//go:embed templates/*
var FS embed.FS

func GetDataForIndividualPersona(personaString string) (user models.User, err error) {
	user, ok := models.UserMap[personaString]
	if !ok {
		return models.User{}, errors.New("user could not be found")
	}
	return user, nil
}

func (t *Template) ExecTemplate(w http.ResponseWriter, baseTemplate string, data any) {
	err := t.htmlTpl.ExecuteTemplate(w, baseTemplate, data)
	if err != nil {
		log.Printf("parsing template: %v", err)
		http.Error(w,
			"There was an error parsing the template",
			http.StatusInternalServerError)
		return
	}
	return
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
