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

type GalleryData struct {
	GalleryFunction string
	TitleInput      inputHTMLAttribs
}

func InitNewGalleryData() GalleryData {
	return GalleryData{
		"new",
		setInputHTMLAttribs(setInputHTMlAttribsValues{"", false, true}),
	}
}

func InitEditGalleryData(loadValue string) GalleryData {
	return GalleryData{
		"edit",
		setInputHTMLAttribs(setInputHTMlAttribsValues{loadValue, false, true}),
	}
}

func InitViewGalleryData(loadValue string) GalleryData {
	return GalleryData{
		"view",
		setInputHTMLAttribs(setInputHTMlAttribsValues{loadValue, true, false}),
	}
}

type setInputHTMlAttribsValues struct {
	loadValue     string
	isSetDisabled bool
	isSetRequired bool
}

func setInputHTMLAttribs(s setInputHTMlAttribsValues) inputHTMLAttribs {
	titleInputHTMLAttribs := inputHTMLAttribs{}
	titleInputHTMLAttribs.SetName("title")
	titleInputHTMLAttribs.SetId("title")
	titleInputHTMLAttribs.SetInputType("text")
	titleInputHTMLAttribs.SetPlaceHolder("Gallery Title")
	titleInputHTMLAttribs.SetLabelText("Title")
	if s.loadValue != "" {
		titleInputHTMLAttribs.SetValue(s.loadValue)
	}
	if s.isSetDisabled {
		titleInputHTMLAttribs.SetIsDisabled()
	}
	if s.isSetRequired {
		titleInputHTMLAttribs.SetIsRequired()
	}
	return titleInputHTMLAttribs
}
