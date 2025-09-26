package views

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"

	"github.com/sohWenMing/lenslocked/models"
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
	UserId           int
	GalleryId        int
	OtherGalleryData any
}

func (g *GalleryData) String() string {
	jsonBytes, _ := json.MarshalIndent(g, "", "    ")
	return string(jsonBytes)
}

func InitNewGalleryData(userId int, loadValue string) GalleryData {
	galleryData := GalleryData{
		UserId:           userId,
		GalleryId:        0,
		OtherGalleryData: InitEditGalleryFunctionAndInputData(loadValue),
	}
	fmt.Println("TOREMOVE: galleryData: ", galleryData.String())
	return galleryData
}

func InitViewGalleryData(userId int, galleryId int, galleryTitle string, exts []string) (GalleryData, error) {
	filePaths, err := getImageFilePaths(galleryId, exts)
	if err != nil {
		return GalleryData{}, err
	}
	galleryData := GalleryData{
		UserId:    userId,
		GalleryId: galleryId,
		OtherGalleryData: struct {
			Title     string
			ImageUrls []string
		}{
			galleryTitle,
			filePaths,
		},
	}
	return galleryData, nil
}

func InitEditGalleryData(userId int, galleryId int, loadTitleValue string, exts []string) (GalleryData, error) {
	filePaths, err := getImageFilePaths(galleryId, exts)
	if err != nil {
		return GalleryData{}, err
	}
	galleryData :=
		GalleryData{
			UserId:    userId,
			GalleryId: galleryId,
			OtherGalleryData: struct {
				ImageUrls []string
				InputData GalleryFunctionToInputData
			}{
				filePaths,
				InitEditGalleryFunctionAndInputData(loadTitleValue),
			},
		}
	fmt.Println("TOREMOVE:  galleryData: ", galleryData.String())
	return galleryData, nil
}

func getImageFilePaths(galleryId int, exts []string) ([]string, error) {
	galleryImages, err := models.GetImagesByGalleryId(galleryId, exts)
	if err != nil {
		return []string{}, err
	}
	filePaths := make([]string, len(galleryImages))
	for i, galleryImage := range galleryImages {
		filePaths[i] = galleryImage.GetPath()
	}
	return filePaths, nil

}

type GalleryFunctionToInputData struct {
	GalleryFunction string
	TitleInput      inputHTMLAttribs
}

func InitNewGalleryFunctionAndInputData() GalleryFunctionToInputData {
	return GalleryFunctionToInputData{
		"new",
		setInputHTMLAttribs(setInputHTMlAttribsValues{"", false, true}),
	}
}

func InitEditGalleryFunctionAndInputData(loadValue string) GalleryFunctionToInputData {
	return GalleryFunctionToInputData{
		"edit",
		setInputHTMLAttribs(setInputHTMlAttribsValues{loadValue, false, true}),
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
