package main

import (
	"embed"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sohWenMing/lenslocked/controllers"
	"github.com/sohWenMing/lenslocked/models"
	"github.com/sohWenMing/lenslocked/views"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func main() {
	envVars, err := models.LoadEnv(".env")
	if err != nil {
		log.Fatal(err)
	}
	dbc, err := models.InitDBConnections()
	if err != nil {
		log.Fatal(err)
	}

	err = models.Migrate(dbc.DB, "migrations", embedMigrations)
	fmt.Println("running migrations on startup")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Migrations successfully ran")

	defer dbc.DB.Close()

	template := views.LoadTemplates()
	//panic would occur if error occured during the loading of templates.

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	userContext := controllers.NewUserContext(dbc.UserService)
	templateHandler := controllers.InitTemplateHandler(template, userContext)

	// ##### Get Method Handlers #####
	// these are not protected routes, so we just use the CookieAuthMiddleWare to test for existence of logged in user
	r.Route("/", func(sr chi.Router) {
		sr.Use(controllers.CookieAuthMiddleWare(dbc.SessionService, nil, false, false))
		sr.Get("/contact", templateHandler("contact.gohtml"))
		sr.Get("/signup", templateHandler("signup.gohtml"))
		sr.Get("/signin", templateHandler("signin.gohtml"))
		sr.Get("/faq", templateHandler("faq.gohtml"))
		sr.Get("/check_email", templateHandler("check_email.gohtml"))
		sr.Get("/", templateHandler("home.gohtml"))
		sr.Get("/forgot_password", templateHandler("forgot_password.gohtml"))
		sr.Get("/reset_password", controllers.TestResetPasswordHandler)
	})

	// these are protected routes, so we use the CookieAuthMiddleWare to test for existence of logged in user and redirect
	// to login if necessary
	r.Route("/user", func(sr chi.Router) {
		sr.Use(controllers.CookieAuthMiddleWare(dbc.SessionService, nil, true, false))
		sr.Use(userContext.SetUserMW())
		sr.Get("/about", templateHandler("user_info.gohtml"))
	})

	r.Get("/test_cookie", templateHandler("test_cookie.gohtml"))
	r.Get("/send_cookie", controllers.TestSendCookie)

	// ##### POST Method Handlers #####
	r.Post("/signup", controllers.HandleSignupForm(dbc))
	r.Post("/signin", controllers.HandleSignInForm(dbc))
	r.Post("/signout", controllers.HandlerSignOut(dbc.SessionService, nil))
	r.Post("/reset_password", controllers.HandleForgotPasswordForm(dbc, envVars.BaseUrl))

	// ##### Not Found Handler #####
	r.NotFound(controllers.ErrNotFoundHandler)

	CSRFMw := controllers.CSRFProtect(envVars.IsDev, envVars.CSRFSecretKey)

	fmt.Println("Starting the server on :3000...")
	log.Fatal(http.ListenAndServe(":3000", CSRFMw(r)))
}
