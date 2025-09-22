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

	mainPagesTemplate := views.LoadPageTemplates(views.MainPagesFS, "templates")

	galleries := &controllers.Galleries{}
	galleries.GalleryService = dbc.GalleryService
	galleries.ConstructNewTemplate(
		&views.GalleryTemplateConstructor{},
		views.GalleryFS,
		[]string{"tailwind_widgets.gohtml",
			"galleries/new_gallery.gohtml",
		},
		"templates")
	galleries.ConstructEditTemplate(
		&views.GalleryTemplateConstructor{},
		views.GalleryFS,
		[]string{"tailwind_widgets.gohtml",
			"galleries/view_edit_gallery.gohtml",
		},
		"templates")
	galleries.ConstructListTemplate(
		&views.GalleryTemplateConstructor{},
		views.GalleryFS,
		[]string{"tailwind_widgets.gohtml",
			"galleries/gallery_index.gohtml",
		},
		"templates")
	//panic would occur if error occured during the loading of templates.

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	userContext := controllers.NewUserContext(dbc.UserService)
	makeHandler, render := controllers.InitTemplateHandler(mainPagesTemplate, userContext)

	// ##### Get Method Handlers #####
	// these are not protected routes, so we just use the CookieAuthMiddleWare to test for existence of logged in user
	r.Route("/", func(sr chi.Router) {
		sr.Use(controllers.CookieAuthMiddleWare(dbc.SessionService, nil, false, false))
		sr.Get("/contact", makeHandler("contact.gohtml"))
		sr.Get("/signup", makeHandler("signup.gohtml"))
		sr.Get("/signin", makeHandler("signin.gohtml"))
		sr.Get("/faq", makeHandler("faq.gohtml"))
		sr.Get("/check_email", makeHandler("check_email.gohtml"))
		sr.Get("/", makeHandler("home.gohtml"))
		sr.Get("/forgot_password", makeHandler("forgot_password.gohtml"))
		sr.Get("/reset_password", makeHandler("reset_password.gohtml"))
		sr.Post("/signup", controllers.HandleSignupForm(dbc, render))
		sr.Post("/signin", controllers.HandleSignInForm(dbc, render))
		sr.Post("/signout", controllers.HandlerSignOut(dbc.SessionService, nil))
		sr.Post("/reset_password", controllers.HandleForgotPasswordForm(dbc, baseUrl, emailService, render))
		sr.Post("/reset_password_submit", controllers.HandlerResetPasswordForm(dbc, render))
	})

	// these are protected routes, so we use the CookieAuthMiddleWare to test for existence of logged in user and redirect
	// to login if necessary
	r.Route("/user", func(sr chi.Router) {
		sr.Use(controllers.CookieAuthMiddleWare(dbc.SessionService, nil, true, false))
		sr.Use(userContext.SetUserMW())
		sr.Get("/about", makeHandler("user_info.gohtml"))
	})
	r.Route("/galleries", func(sr chi.Router) {
		sr.Get("/{id}", galleries.View(dbc.GalleryService))
		sr.Group(func(sr chi.Router) {
			sr.Use(controllers.CookieAuthMiddleWare(dbc.SessionService, nil, true, false))
			sr.Use(userContext.SetUserMW())
			sr.Get("/new_gallery", galleries.New)
			sr.Get("/{id}/edit", galleries.Edit(dbc.GalleryService))
			sr.Get("/list", galleries.List)
			sr.Post("/new", galleries.Create)
			sr.Post("/edit", galleries.HandleEdit(dbc.GalleryService))
			sr.Post("/{id}/delete", galleries.HandleDelete(dbc.GalleryService))
		})
	})

	r.Get("/test_cookie", makeHandler("test_cookie.gohtml"))
	r.Get("/send_cookie", controllers.TestSendCookie)
	r.Get("/test_alert", makeHandler("test_alert.gohtml"))

	// ##### POST Method Handlers #####

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
