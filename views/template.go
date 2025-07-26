package views

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type Template struct {
	htmlTpl *template.Template
}

type TplMap map[string]Template

var tplStrings = []string{
	"home.gohtml",
	"contact.gohtml",
	// "faq.gohtml",
	// "persona.gohtml",
	"layout-parts.gohtml",
}

var BaseTemplateToData = map[string]any{
	"home.gohtml":    nil,
	"contact.gohtml": nil,
}

//go:embed templates/*
var FS embed.FS

func (t *Template) ExecTemplate(w http.ResponseWriter, baseTemplate string) (err error) {
	data := BaseTemplateToData[baseTemplate]
	w.Header().Set("content-type", "text/html; charset=utf-8")
	err = t.htmlTpl.ExecuteTemplate(w, baseTemplate, data)
	if err != nil {
		log.Printf("parsing template: %v", err)
		http.Error(w,
			"There was an error parsing the template",
			http.StatusInternalServerError)
		return
	}
	return nil
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
