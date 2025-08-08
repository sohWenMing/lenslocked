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

/*
The Template struct is used to house the type *template.Template so that a the method ExecTemplate can be
attached to it.

ExecTemplateWithCSTF - allows us to pass in the csrfField, which will in turn be passed on to the function defined in
cloned.Funcs

# ExecTemplate - normal execution of template with the need fo csrfField

In both cases, we need to clone the template because the one that is for CSRF requires cloning to be safe, due to the
function mutating the template at each request
*/
func (t *Template) ExecTemplateWithCSRF(
	w http.ResponseWriter,
	r *http.Request,
	csrfField template.HTML,
	baseTemplate string,
	data any) {
	cloned, err := t.htmlTpl.Clone()
	if err != nil {
		log.Printf("parsing template: %v", err)
		http.Error(w,
			"There was an error parsing the template",
			http.StatusInternalServerError)
		return
	}
	cloned = cloned.Funcs(
		template.FuncMap{
			"csrfField": func() template.HTML {
				return csrfField
			},
		},
	)
	err = cloned.ExecuteTemplate(w, baseTemplate, data)
	if err != nil {
		log.Printf("parsing template: %v", err)
		http.Error(w,
			"There was an error parsing the template",
			http.StatusInternalServerError)
		return
	}
}

func (t *Template) ExecTemplate(w http.ResponseWriter, r *http.Request, baseTemplate string, data any) {
	clone, err := t.htmlTpl.Clone()
	if err != nil {
		log.Printf("parsing template: %v", err)
		http.Error(w,
			"There was an error parsing the template",
			http.StatusInternalServerError)
		return
	}

	err = clone.ExecuteTemplate(w, baseTemplate, data)
	if err != nil {
		log.Printf("parsing template: %v", err)
		http.Error(w,
			"There was an error parsing the template",
			http.StatusInternalServerError)
		return
	}
}

var tplStrings = []string{
	"home.gohtml",
	"contact.gohtml",
	"faq.gohtml",
	"persona.gohtml",
	"persona_multiple.gohtml",
	"tailwind_widgets.gohtml",
	"signup.gohtml",
	"signin.gohtml",
	"practice_form.gohtml",
	"test_cookie.gohtml",
}

// list of all the defined templates that will be used for rendering

var BaseTemplateToData = map[string]any{
	"home.gohtml":             nil,
	"contact.gohtml":          nil,
	"faq.gohtml":              models.QuestionsToAnswers,
	"signup.gohtml":           SignUpSignInFormData,
	"signin.gohtml":           SignUpSignInFormData,
	"practice_form.gohtml":    nil,
	"test_cookie.gohtml":      nil,
	"persona_multiple.gohtml": nil,
}

// defines the data that will be passed in at execution time for each base template

//go:embed templates/*
var FS embed.FS

func LoadTemplates() (tpl *Template) {
	tpl = &Template{}
	loadedTemplate := template.New("base")
	//sets up basically an empty template, so that we can load functions in to it BEFORE we actually parse all the rest of templates
	loadedTemplate = loadedTemplate.Funcs(
		template.FuncMap{
			"csrfField": func() template.HTML {
				return `<input type="hidden" />`
			},
		},
	)

	//this is a placeholder function - we need this or else
	templateStrings := getTemplatePaths(tplStrings, "templates")
	loadedTemplate = TemplateMust(loadedTemplate.ParseFS(FS, templateStrings...))
	tpl.htmlTpl = loadedTemplate
	return tpl
}

/*
used to load up and parse all templates at the beginning of execution of the program
TemplateMust function will panic any error found during parsing, which will shut down execution of the program
*/

func getTemplatePaths(tplStrings []string, baseFolderName string) []string {
	fullPaths := make([]string, len(tplStrings))
	for i, tplString := range tplStrings {
		fullPath := fmt.Sprintf("%s/%s", baseFolderName, tplString)
		fullPaths[i] = fullPath
	}
	return fullPaths
}

// helper function - used to create the final string slice of template paths that will be parsed

func TemplateMust(t *template.Template, err error) *template.Template {
	if err != nil {
		panic(err)
	}
	return t
}
