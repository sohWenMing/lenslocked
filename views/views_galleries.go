package views

import (
	"embed"
	"html/template"
)

// here i need something, which can be used to comstruct and pass out a views.Template struct
// so the method, the output of the method needs to be a views.Template

//go:embed templates/*
var GalleryFS embed.FS

// Loads all the templates that are defined in template strings, which should be present in the
// embed.FS
type GalleryTemplateConstructor struct{}

func (n *GalleryTemplateConstructor) ConstructTemplate(fs embed.FS, templateStrings []string, baseFolderName string) *Template {
	tpl := &Template{}
	loadedTemplate := template.New("base")
	// all i want to here, is to get all the files, and the parse them to get them ready
	loadedTemplate = loadedTemplate.Funcs(
		template.FuncMap{
			"csrfField": func() template.HTML {
				return `<input type="hidden" />`
			},
			"errors": func() []string {
				return []string{}
			},
		},
	)
	tplStrings := getTemplatePaths(templateStrings, baseFolderName)
	loadedHTMLTemplate := TemplateMust(loadedTemplate.ParseFS(fs, tplStrings...))
	tpl.htmlTpl = loadedHTMLTemplate
	return tpl
}

type NewGalleryData struct {
	FirstInput inputHTMLAttribs
}

func InitNewGalleryData() NewGalleryData {
	firstInputHtmlAttribs := inputHTMLAttribs{}
	firstInputHtmlAttribs.SetName("title")
	firstInputHtmlAttribs.SetId("title")
	firstInputHtmlAttribs.SetInputType("text")
	firstInputHtmlAttribs.SetPlaceHolder("Gallery Title")
	firstInputHtmlAttribs.SetLabelText("Title")
	firstInputHtmlAttribs.SetIsRequired(true)
	return NewGalleryData{
		firstInputHtmlAttribs,
	}
}
