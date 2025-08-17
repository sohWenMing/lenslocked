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

	// ##### Get Method Handlers #####
	r.Get("/contact", controllers.HandlerExecuteTemplate(template, "contact.gohtml"))
	r.Get("/signup", controllers.HandlerExecuteTemplate(template, "signup.gohtml"))
	r.Get("/signin", controllers.HandlerExecuteTemplate(template, "signin.gohtml"))
	r.Get("/faq", controllers.HandlerExecuteTemplate(template, "faq.gohtml"))

	r.Route("/user", func(sr chi.Router) {
		sr.Use(controllers.ProtectedCookieAuthMiddleWare(dbc.SessionService, nil, false))
		sr.Get("/about", controllers.HandlerExecuteTemplate(template, "persona_multiple.gohtml"))
	})
	// these are protected subrroutes, where we would want ot check for the existence of a cookie in the request

	r.Get("/forgot_password", controllers.TestHandler("To do - forgot password page"))
	r.Get("/test_cookie", controllers.HandlerExecuteTemplate(template, "test_cookie.gohtml"))
	r.Get("/send_cookie", controllers.TestSendCookie)

	r.Get("/", controllers.HandlerExecuteTemplate(template, "home.gohtml"))
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
