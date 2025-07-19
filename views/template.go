package views

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

type Template struct {
	htmlTpl *template.Template
}

type TplMap map[string]Template

var tplStrings = []string{
	"home.gohtml",
	"contact.gohtml",
	"faq.gohtml",
}

func LoadTemplates(relPath string) (tplMap *TplMap, err error) {

	workingMap := TplMap{}

	for _, tplString := range tplStrings {
		tplPath := filepath.Join(relPath, tplString)
		tpl, err := template.ParseFiles(tplPath)
		if err != nil {
			return nil, err
		}
		workingTemplate := Template{
			tpl,
		}
		workingMap[tplString] = workingTemplate
	}
	return &workingMap, nil
}

func (t *Template) Execute(w http.ResponseWriter, data interface{}) (err error) {
	w.Header().Set("content-type", "text/html; charset=utf-8")
	err = t.htmlTpl.Execute(w, data)
	if err != nil {
		log.Printf("parsing template: %v", err)
		http.Error(w,
			"There was an error parsing the template",
			http.StatusInternalServerError)
		return
	}
	return nil
}
