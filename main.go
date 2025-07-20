package main

import (
	"fmt"
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
	tplMap := views.LoadTemplates("./templates")
	//panic would occur if error occured during the loading of templates.

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", controllers.GetHandlerFromTplMap(*tplMap, "home.gohtml"))
	r.Get("/contact", controllers.GetHandlerFromTplMap(*tplMap, "contact.gohtml"))
	r.Get("/faq", controllers.GetHandlerFromTplMap(*tplMap, "faq.gohtml"))
	r.Get("/about/{persona}", controllers.GetHandlerforAbout(*tplMap, "persona.gohtml"))
	r.NotFound(controllers.ErrNotFoundHandler)
	fmt.Println("Starting the server on :3000...")
	http.ListenAndServe(":3000", r)
}
