package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sohWenMing/lenslocked/controllers"
	"github.com/sohWenMing/lenslocked/gomailer"
	"github.com/sohWenMing/lenslocked/migrations"
	"github.com/sohWenMing/lenslocked/models"
	"github.com/sohWenMing/lenslocked/services"
	"github.com/sohWenMing/lenslocked/views"
)

type config struct {
	isDev         bool
	baseUrl       string
	csrfSecretKey string
	emailEnvVars  *models.EmailEnvs
	pgConfig      models.PgConfig
}

func loadEnvConfig() (*config, error) {
	envVars, err := loadEnvVars()
	if err != nil {
		return nil, err
	}
	isDev, err := readIsDev(envVars)
	if err != nil {
		return nil, err
	}
	baseUrl, err := readBaseUrl(envVars)
	if err != nil {
		return nil, err
	}
	csrfSecretKey, err := readCSRFSecretKey(envVars)
	if err != nil {
		return nil, err
	}
	emailEnvVars, err := getEmailEnvVars(envVars)
	if err != nil {
		return nil, err
	}
	pgConfig, err := envVars.LoadPgConfig()
	if err != nil {
		return nil, err
	}
	return &config{
		isDev, baseUrl, csrfSecretKey, emailEnvVars, pgConfig,
	}, nil
}

func main() {
	cfg, err := loadEnvConfig()
	if err != nil {
		panic(err)
	}
	err = run(cfg)
	if err != nil {
		panic(err)
	}
}

func run(cfg *config) error {

	initGoMailer := gomailer.NewGoMailer(
		cfg.emailEnvVars.Host,
		cfg.emailEnvVars.Username,
		cfg.emailEnvVars.Password,
		cfg.emailEnvVars.Port)

	emailService := services.InitEmailService(initGoMailer, services.LoadEmailTemplates())

	dbc, err := models.InitDBConnections(cfg.pgConfig)
	if err != nil {
		return err
	}
	err = models.Migrate(dbc.DB, ".", migrations.GetMigrations())
	if err != nil {
		return err
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
			"galleries/edit_gallery.gohtml",
		},
		"templates")
	galleries.ConstructListTemplate(
		&views.GalleryTemplateConstructor{},
		views.GalleryFS,
		[]string{"tailwind_widgets.gohtml",
			"galleries/gallery_index.gohtml",
		},
		"templates")
	galleries.ConstructViewTemplate(
		&views.GalleryTemplateConstructor{},
		views.GalleryFS,
		[]string{"tailwind_widgets.gohtml",
			"galleries/view_gallery.gohtml",
		},
		"templates")
	//panic would occur if error occured during the loading of templates.
	r := chi.NewRouter()
	// r.Use(middleware.Logger)
	// r.Handle("/images/*", http.StripPrefix("/images/", models.LoadImageFileServer("./images")))

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
		sr.Post("/reset_password", controllers.HandleForgotPasswordForm(dbc, cfg.baseUrl, emailService, render))
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
		sr.Group(func(sr chi.Router) {
			sr.Use(controllers.CookieAuthMiddleWare(dbc.SessionService, nil, false, false))
			sr.Use(userContext.SetUserMW())
			sr.Get("/{id}", galleries.View(dbc.GalleryService))
			sr.Handle("/{id}/images/{filename}", controllers.ServeImage())
		})
		sr.Group(func(sr chi.Router) {
			sr.Use(middleware.Logger)
			sr.Use(controllers.CookieAuthMiddleWare(dbc.SessionService, nil, true, false))
			sr.Use(userContext.SetUserMW())
			sr.Get("/new_gallery", galleries.New)
			sr.Get("/{id}/edit", galleries.Edit(dbc.GalleryService))
			sr.Get("/list", galleries.List)
			sr.Post("/new", galleries.Create)
			sr.Post("/edit", galleries.HandleEdit(dbc.GalleryService))
			sr.Post("/{id}/delete", galleries.HandleDelete(dbc.GalleryService))
			sr.Post("/{id}/images/{filename}/delete", galleries.DeleteImage(dbc.GalleryService))
			sr.Post("/{id}/images", galleries.UploadImage(dbc.GalleryService))
		})
	})

	r.Get("/test_cookie", makeHandler("test_cookie.gohtml"))
	r.Get("/send_cookie", controllers.TestSendCookie)
	r.Get("/test_alert", makeHandler("test_alert.gohtml"))

	// ##### POST Method Handlers #####

	// ##### Not Found Handler #####
	r.NotFound(controllers.ErrNotFoundHandler)

	CSRFMw := controllers.CSRFProtect(cfg.isDev, cfg.csrfSecretKey)

	fmt.Println("Starting the server on :3000...")
	return (http.ListenAndServe(":3000", CSRFMw(r)))
}

func readBaseUrl(envVars *models.Envs) (string, error) {
	envBaseUrl, err := envVars.GetBaseURL()
	if err != nil {
		return "", err
	}
	return envBaseUrl, nil
}

func readCSRFSecretKey(envVars *models.Envs) (string, error) {
	envcsrfSecretKey, err := envVars.GetCSRFSecretKey()
	if err != nil {
		return "", err
	}
	return envcsrfSecretKey, nil
}

func readIsDev(envVars *models.Envs) (bool, error) {
	envIsDev, err := envVars.GetIsDev()
	if err != nil {
		return false, err
	}
	return envIsDev, nil
}

func getEmailEnvVars(envVars *models.Envs) (*models.EmailEnvs, error) {
	emailEnvs, err := envVars.LoadEmailEnvs()
	if err != nil {
		return nil, err
	}
	return emailEnvs, nil
}

func loadEnvVars() (*models.Envs, error) {
	envVars, err := models.LoadEnv(".env")
	if err != nil {
		return nil, err
	}
	return envVars, nil
}
