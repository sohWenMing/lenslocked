package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sohWenMing/lenslocked/controllers"
	"github.com/sohWenMing/lenslocked/views"
)

type templateDirAndPath struct {
	directory string
	path      string
}

func main() {
	template := views.LoadTemplates()
	//panic would occur if error occured during the loading of templates.

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", controllers.HandlerExecuteTemplate(*template, "home.gohtml", nil))
	r.Get("/contact", controllers.HandlerExecuteTemplate(*template, "contact.gohtml", nil))
	r.Get("/faq", controllers.HandlerExecuteTemplate(*template, "faq.gohtml",
		views.BaseTemplateToData["faq.gohtml"]))
	r.Get("/about", controllers.HandlerExecuteTemplate(*template, "persona_multiple.gohtml",
		views.BaseTemplateToData["persona_multiple.gohtml"]))
	r.Get("/about/{persona}", controllers.HandlerForIndividualUser(*template))
	r.NotFound(controllers.ErrNotFoundHandler)
	fmt.Println("Starting the server on :3000...")
	log.Fatal(http.ListenAndServe(":3000", r))
}
