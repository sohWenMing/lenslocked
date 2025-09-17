package controllers

import (
	"embed"
	"fmt"
	"net/http"

	"github.com/sohWenMing/lenslocked/models"
	"github.com/sohWenMing/lenslocked/views"
)

/*
Galleries holds all the templates that have been parsed and made ready for the use of generating gallery related templates,
In addition to the GalleryService which provides connections to gallery related database operations.
*/
type Galleries struct {
	Templates struct {
		New *views.Template
	}
	GalleryService *models.GalleryService
}

// constructor function used to initialise the New template to Galleries struct
func (g *Galleries) ConstructNewTemplate(constructor GalleryTemplateConstructor, fs embed.FS, templateStrings []string, baseFolderName string) {
	g.Templates.New = constructor.ConstructTemplate(fs, templateStrings, baseFolderName)
}

type GalleryTemplateConstructor interface {
	ConstructTemplate(fs embed.FS, templateStrings []string, baseFolderName string) *views.Template
}

/*
New is used render a form that for the creation of of a new gallery.
Will be used for the first time render of the form, so there is no error handling in place
*/
func (g *Galleries) New(w http.ResponseWriter, r *http.Request) {
	// need to get in the user context, after getting in the information, we want to always be able
	csrfToken := GetCSRFTokenFromRequest(r)
	userId, _ := GetUserIdFromRequestContext(r)
	g.Templates.New.ExecTemplateWithCSRF(w, r, csrfToken, "new_gallery.gohtml", initNewGalleryData(userId), nil)
}

func (g *Galleries) Create(w http.ResponseWriter, r *http.Request) {
	csrfToken := GetCSRFTokenFromRequest(r)
	userId, isFound := GetUserIdFromRequestContext(r)
	if !isFound {
		http.Redirect(w, r, "/signin", http.StatusBadRequest)
		return
	}
	title := r.FormValue("title")
	if title == "" {
		g.Templates.New.ExecTemplateWithCSRF(w, r, csrfToken, "new_gallery.gohtml", initNewGalleryData(userId), []string{"mandatory inputs were not filled"})
		return
	}
	w.WriteHeader(http.StatusOK)
	gallery, err := g.GalleryService.Create(title, userId)
	if err != nil {
		g.Templates.New.ExecTemplateWithCSRF(w, r, csrfToken, "new_gallery.gohtml", initNewGalleryData(userId), []string{err.Error()})
		return
	}
	editPath := fmt.Sprintf("/galleries/%d/edit", gallery.ID)
	fmt.Println("editPath: ", editPath)
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

type NewGalleryData struct {
	UserId         int
	NewGalleryData views.NewGalleryData
}

func initNewGalleryData(userId int) NewGalleryData {
	return NewGalleryData{
		UserId:         userId,
		NewGalleryData: views.InitNewGalleryData(),
	}

}

//we want to create a value that:
// checks that the the user id, and title is submitted
