package main

import (
	"embed"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sohWenMing/lenslocked/controllers"
	"github.com/sohWenMing/lenslocked/gomailer"
	"github.com/sohWenMing/lenslocked/models"
	"github.com/sohWenMing/lenslocked/services"
	"github.com/sohWenMing/lenslocked/views"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

var isDev bool
var baseUrl string
var csrfSecretKey string

func main() {
	envVars := loadEnvVars()

	emailEnvVars := getEmailEnvVars(envVars)
	initGoMailer := gomailer.NewGoMailer(emailEnvVars.Host, emailEnvVars.Username, emailEnvVars.Password, emailEnvVars.Port)
	emailService := services.InitEmailService(initGoMailer, services.LoadEmailTemplates())

	setIsDev(envVars)
	setCSRFSecretKey(envVars)
	setBaseUrl(envVars)

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

	template := views.LoadPageTemplates()
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
		sr.Get("/reset_password", templateHandler("reset_password.gohtml"))
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
	r.Post("/reset_password", controllers.HandleForgotPasswordForm(dbc, baseUrl, emailService))
	r.Post("/reset_password_submit", controllers.HandlerResetPasswordForm(dbc))

	// ##### Not Found Handler #####
	r.NotFound(controllers.ErrNotFoundHandler)

	CSRFMw := controllers.CSRFProtect(isDev, csrfSecretKey)

	fmt.Println("Starting the server on :3000...")
	log.Fatal(http.ListenAndServe(":3000", CSRFMw(r)))
}

func setBaseUrl(envVars *models.Envs) {
	envBaseUrl, err := envVars.GetBaseURL()
	if err != nil {
		log.Fatal(err)
	}
	baseUrl = envBaseUrl
}

func setCSRFSecretKey(envVars *models.Envs) {
	envcsrfSecretKey, err := envVars.GetCSRFSecretKey()
	if err != nil {
		log.Fatal(err)
	}
	csrfSecretKey = envcsrfSecretKey
}

func setIsDev(envVars *models.Envs) {
	envIsDev, err := envVars.GetIsDev()
	if err != nil {
		log.Fatal(err)
	}
	isDev = envIsDev
}

func getEmailEnvVars(envVars *models.Envs) *models.EmailEnvs {
	emailEnvs, err := envVars.LoadEmailEnvs()
	if err != nil {
		log.Fatal(err)
	}
	return emailEnvs
}

func loadEnvVars() *models.Envs {
	envVars, err := models.LoadEnv(".env")
	if err != nil {
		log.Fatal(err)
	}
	return envVars
}
