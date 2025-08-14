package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sohWenMing/lenslocked/controllers"
	"github.com/sohWenMing/lenslocked/models"
	"github.com/sohWenMing/lenslocked/views"
)

type templateDirAndPath struct {
	directory string
	path      string
}

func main() {
	envVars, err := models.LoadEnv(".env")
	dbc, err := models.InitDBConnections()
	if err != nil {
		log.Fatal(err)
	}
	defer dbc.DB.Close()

	template := views.LoadTemplates()
	formNameToLoader := controllers.InitFormNameToLoader(template)
	//panic would occur if error occured during the loading of templates.

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// ##### Get Method Handlers #####
	r.Get("/contact", controllers.HandlerExecuteTemplate(template, "contact.gohtml", nil))
	r.Get("/signup", formNameToLoader["signup_form"].Load)
	r.Get("/signin", formNameToLoader["signin_form"].Load)
	r.Get("/faq", controllers.HandlerExecuteTemplate(template, "faq.gohtml",
		views.BaseTemplateToData["faq.gohtml"]))

	r.Route("/user", func(sr chi.Router) {
		sr.Use(controllers.CookieAuthMiddleWare(dbc.SessionService, nil, time.Now()))
		sr.Get("/about", controllers.HandlerExecuteTemplate(template, "persona_multiple.gohtml",
			views.BaseTemplateToData["persona_multiple.gohtml"]))
	})
	// these are protected subrroutes, where we would want ot check for the existence of a cookie in the request

	r.Get("/forgot_password", controllers.TestHandler("To do - forgot password page"))
	r.Get("/test_cookie", controllers.HandlerExecuteTemplate(template, "test_cookie.gohtml", views.BaseTemplateToData["test_cookie.gohtml"]))
	r.Get("/send_cookie", controllers.TestSendCookie)
	r.Get("/signout", controllers.ProcessSignOut(dbc.SessionService, nil))

	r.Get("/", controllers.HandlerExecuteTemplate(template, "home.gohtml", nil))
	// ##### POST Method Handlers #####
	r.Post("/signup", controllers.HandleSignupForm(dbc))
	r.Post("/signin", controllers.HandleSignInForm(dbc))

	// ##### Not Found Handler #####
	r.NotFound(controllers.ErrNotFoundHandler)

	CSRFMw := controllers.CSRFProtect(envVars.IsDev, envVars.CSRFSecretKey)

	fmt.Println("Starting the server on :3000...")
	log.Fatal(http.ListenAndServe(":3000", CSRFMw(r)))
}
