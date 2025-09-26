package controllers

import (
	"embed"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/sohWenMing/lenslocked/models"
	"github.com/sohWenMing/lenslocked/views"
)

/*
Galleries holds all the templates that have been parsed and made ready for the use of generating gallery related templates,
In addition to the GalleryService which provides connections to gallery related database operations.
*/
type Galleries struct {
	Templates struct {
		New  *views.Template
		Edit *views.Template
		List *views.Template
		View *views.Template
	}
	GalleryService *models.GalleryService
}

// constructor function used to initialise the New template to Galleries struct
func (g *Galleries) ConstructNewTemplate(constructor GalleryTemplateConstructor, fs embed.FS, templateStrings []string, baseFolderName string) {
	g.Templates.New = constructor.ConstructTemplate(fs, templateStrings, baseFolderName)
}
func (g *Galleries) ConstructEditTemplate(constructor GalleryTemplateConstructor, fs embed.FS, templateStrings []string, baseFolderName string) {
	g.Templates.Edit = constructor.ConstructTemplate(fs, templateStrings, baseFolderName)
}
func (g *Galleries) ConstructListTemplate(constructor GalleryTemplateConstructor, fs embed.FS, templateStrings []string, baseFolderName string) {
	g.Templates.List = constructor.ConstructTemplate(fs, templateStrings, baseFolderName)
}
func (g *Galleries) ConstructViewTemplate(constructor GalleryTemplateConstructor, fs embed.FS, templateStrings []string, baseFolderName string) {
	g.Templates.View = constructor.ConstructTemplate(fs, templateStrings, baseFolderName)
}

type GalleryTemplateConstructor interface {
	ConstructTemplate(fs embed.FS, templateStrings []string, baseFolderName string) *views.Template
}

type GalleryListing struct {
	Id    int
	Title string
}
type GalleryListData struct {
	UserId          int
	GalleryListings []GalleryListing
}

func (g *Galleries) List(w http.ResponseWriter, r *http.Request) {
	userId, _ := GetUserIdFromRequestContext(r)
	csrfToken := GetCSRFTokenFromRequest(r)
	galleries, err := g.GalleryService.GetGalleryListByUserId(userId)
	fmt.Println("galleries: ", galleries)
	if err != nil {
		http.Error(w, fmt.Sprintf("error %s", err.Error()), http.StatusInternalServerError)
		return
	}
	galleryListings := make([]GalleryListing, len(galleries))
	for i, gallery := range galleries {
		galleryListings[i] = GalleryListing{
			Id:    gallery.ID,
			Title: gallery.Title,
		}
	}
	galleryData := GalleryListData{
		userId, galleryListings,
	}
	g.Templates.List.ExecTemplateWithCSRF(w, r, csrfToken, "gallery_index.gohtml", galleryData, nil)
}

/*
New is used render a form that for the creation of of a new gallery.
Will be used for the first time render of the form, so there is no error handling in place
*/
func (g *Galleries) New(w http.ResponseWriter, r *http.Request) {
	// need to get in the user context, after getting in the information, we want to always be able
	csrfToken := GetCSRFTokenFromRequest(r)
	userId, _ := GetUserIdFromRequestContext(r)
	g.Templates.New.ExecTemplateWithCSRF(w, r, csrfToken, "new_gallery.gohtml", views.InitNewGalleryData(userId, ""), nil)
}

/*
Create is the handler used when a user attempts to send a Post request to create a new gallery.
If there is an error, will render the same create template that was initially loaded, but with the error
else, will process request and then redirect to the edit page
*/
func (g *Galleries) Create(w http.ResponseWriter, r *http.Request) {
	csrfToken := GetCSRFTokenFromRequest(r)
	userId, _ := GetUserIdFromRequestContext(r)
	title := r.FormValue("title")
	if title == "" {
		g.Templates.New.ExecTemplateWithCSRF(w, r, csrfToken, "new_gallery.gohtml", views.InitNewGalleryData(userId, title), []string{"mandatory inputs were not filled"})
		return
	}
	gallery, err := g.GalleryService.Create(title, userId)
	if err != nil {
		g.Templates.New.ExecTemplateWithCSRF(w, r, csrfToken, "new_gallery.gohtml", views.InitNewGalleryData(userId, title), []string{err.Error()})
		return
	}
	http.Redirect(w, r, getEditPath(gallery.ID), http.StatusFound)
}

func (g *Galleries) Edit(gs *models.GalleryService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		csrfToken := GetCSRFTokenFromRequest(r)
		userId, _ := GetUserIdFromRequestContext(r)
		gallery, err := getValidatedUserGallery(r, gs, userId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		galleryData, err := views.InitEditGalleryData(userId, gallery.ID, gallery.Title, g.GalleryService.GetImageExtensions())
		if err != nil {
			http.Error(w, "Internal SErver Error", http.StatusInternalServerError)
			return
		}
		g.Templates.Edit.ExecTemplateWithCSRF(w, r, csrfToken, "edit_gallery.gohtml", galleryData, nil)
	}
}

func getValidatedUserGallery(r *http.Request, gs *models.GalleryService, userId int) (gallery *models.Gallery, err error) {
	gallery, err = getGalleryByRequestGalleryId(r, gs)
	if err != nil {
		return nil, err
	}
	if userId != gallery.UserID {
		return nil, errors.New("user is not owner of gallery")
	}
	return gallery, nil
}
func (g *Galleries) View(gs *models.GalleryService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		csrfToken := GetCSRFTokenFromRequest(r)
		gallery, err := getGalleryByRequestGalleryId(r, gs)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		userId, _ := GetUserIdFromRequestContext(r)
		galleryData, err := views.InitViewGalleryData(userId, gallery.ID, gallery.Title, g.GalleryService.GetImageExtensions())

		fmt.Println("TOREMOVE: galleryData: ", galleryData.String())
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		g.Templates.View.ExecTemplateWithCSRF(w, r, csrfToken, "view_gallery.gohtml", galleryData, nil)
	}
}

func (g *Galleries) DeleteImage(gs *models.GalleryService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		filename := chi.URLParam(r, "filename")
		fmt.Println("filename: ", filename)
		userId, _ := GetUserIdFromRequestContext(r)
		gallery, err := getValidatedUserGallery(r, gs, userId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		evalpath := fmt.Sprintf("./images/%d/%s", gallery.ID, filename)
		err = os.Remove(evalpath)
		if err != nil {
			http.Error(w, fmt.Sprintf("error occured when trying to delete: %s", err.Error()), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, fmt.Sprintf("/galleries/%d/edit", gallery.ID), http.StatusFound)

	}
}

func ServeImage() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			galleryIdString := chi.URLParam(r, "id")
			galleryId, err := strconv.Atoi(galleryIdString)
			if err != nil {
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			}
			fileName := chi.URLParam(r, "filename")

			filePath := fmt.Sprintf("./images/%d/%s", galleryId, fileName)
			if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
				http.Error(w, "file does not exist", http.StatusNotFound)
				return
			} else if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			} else {
				http.ServeFile(w, r, filePath)
			}
		})
}

func (g *Galleries) HandleEdit(gs *models.GalleryService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		userId, _ := GetUserIdFromRequestContext(r)

		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Internal Error", http.StatusInternalServerError)
			return
		}
		id := r.Form.Get("gallery-id")
		galleryId, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		title := r.Form.Get("title")
		if title == "" {
			http.Error(w, "Mandatory information not filled", http.StatusBadRequest)
			return
		}
		gallery, err := gs.GetById(galleryId)
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		if gallery.UserID != userId {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		err = gs.UpdateTitle(galleryId, title)
		if err != nil {
			http.Error(w, "Internal Error", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/galleries/list", http.StatusFound)
	}
}
func (g *Galleries) HandleDelete(gs *models.GalleryService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		userId, _ := GetUserIdFromRequestContext(r)

		gallery, err := getGalleryByRequestGalleryId(r, gs)
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		if gallery.UserID != userId {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		err = gs.DeleteById(gallery.ID)
		if err != nil {
			http.Error(w, "Internal Error", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/galleries/list", http.StatusFound)
	}
}

func getGalleryByRequestGalleryId(r *http.Request, galleryService *models.GalleryService) (gallery *models.Gallery, err error) {
	galleryId, err := getGalleryIdFromRequest(r)
	if err != nil {
		fmt.Println("err: ", err.Error())
		return nil, err
	}
	gallery, err = getGalleryById(galleryId, galleryService)
	if err != nil {
		return nil, err
	}
	return gallery, nil
}

func getGalleryIdFromRequest(r *http.Request) (galleryId int, err error) {
	id := chi.URLParam(r, "id")
	galleryId, err = strconv.Atoi(id)
	if err != nil {
		return 0, err
	}
	return galleryId, nil
}

func getGalleryById(galleryId int, galleryService *models.GalleryService) (gallery *models.Gallery, err error) {
	gallery, err = galleryService.GetById(galleryId)
	if err != nil {
		return nil, err
	}
	return gallery, nil
}

func getEditPath(galleryId int) string {
	editPath := fmt.Sprintf("/galleries/%d/edit", galleryId)
	return editPath
}
