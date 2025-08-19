package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sohWenMing/lenslocked/controllers"
	"github.com/sohWenMing/lenslocked/models"
	"github.com/sohWenMing/lenslocked/views"
)

func main() {
	envVars, err := models.LoadEnv(".env")
	if err != nil {
		log.Fatal(err)
	}
	dbc, err := models.InitDBConnections()
	if err != nil {
		log.Fatal(err)
	}
	defer dbc.DB.Close()

	template := views.LoadTemplates()
	//panic would occur if error occured during the loading of templates.

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	handlerExecuteTemplateFunc := controllers.InitHandlerExecuteTemplateFunc(template, dbc.UserService)

	// ##### Get Method Handlers #####
	// these are not protected routes, so we just use the CookieAuthMiddleWare to test for existence of logged in user
	r.Route("/", func(sr chi.Router) {
		sr.Use(controllers.CookieAuthMiddleWare(dbc.SessionService, nil, false, false))
		sr.Get("/contact", handlerExecuteTemplateFunc("contact.gohtml"))
		sr.Get("/signup", handlerExecuteTemplateFunc("signup.gohtml"))
		sr.Get("/signin", handlerExecuteTemplateFunc("signin.gohtml"))
		sr.Get("/faq", handlerExecuteTemplateFunc("faq.gohtml"))
		sr.Get("/", handlerExecuteTemplateFunc("home.gohtml"))
	})

	// these are protected routes, so we use the CookieAuthMiddleWare to test for existence of logged in user and redirect
	// to login if necessary
	r.Route("/user", func(sr chi.Router) {
		sr.Use(controllers.CookieAuthMiddleWare(dbc.SessionService, nil, true, false))
		sr.Get("/about", handlerExecuteTemplateFunc("persona_multiple.gohtml"))
	})

	r.Get("/forgot_password", controllers.TestHandler("To do - forgot password page"))
	r.Get("/test_cookie", handlerExecuteTemplateFunc("test_cookie.gohtml"))
	r.Get("/send_cookie", controllers.TestSendCookie)

	// ##### POST Method Handlers #####
	r.Post("/signup", controllers.HandleSignupForm(dbc))
	r.Post("/signin", controllers.HandleSignInForm(dbc))
	r.Post("/signout", controllers.HandlerSignOut(dbc.SessionService, nil))

	// ##### Not Found Handler #####
	r.NotFound(controllers.ErrNotFoundHandler)

	CSRFMw := controllers.CSRFProtect(envVars.IsDev, envVars.CSRFSecretKey)

	fmt.Println("Starting the server on :3000...")
	log.Fatal(http.ListenAndServe(":3000", CSRFMw(r)))
}
