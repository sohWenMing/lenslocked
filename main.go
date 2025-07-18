package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const faq string = `
	<h1>FAQ Page</h1>
	<ul>
		<li>
			<b>Is there a free version?</b>
			<br>
			Yes! We offer a free trial for 30 days on any paid plans
		</li>
		<li>
			<b>What are the support hours?</b>
			<br>
			We have support staff answering emails 24/7, though response times might be a bit slower on weekends
		</li>
		<li>
			<b>How do I contact support staff?</b>
			<br>
			Email us - <a href="mailto:support@lenslocked.com">support@lenslocked.com</a>
		</li>
	</ul>
`

type templateDirAndPath struct {
	directory string
	path      string
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/html")
	executeTemplate(&w, templateDirAndPath{"templates", "home.gohtml"}, nil)
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("contene-type", "text/html; charset=utf-8")
	executeTemplate(&w, templateDirAndPath{"templates", "contact.gohtml"}, nil)
}
func errNotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "404 not found")
}
func faqHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/html; charset=utf-8")
	fmt.Fprintf(w, faq)
}
func checkParamHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/html; charset=utf-8")
	param := chi.URLParam(r, "param")
	html := fmt.Sprintf("<h1>%s</h1>", param)
	fmt.Fprintf(w, html)
}

func bioHandler(w http.ResponseWriter, r *http.Request) {
	bio := `&lt;script&gt;alert(&quot;Hi!&quot;);&lt;/script&gt;`
	w.Header().Set("content-type", "text/html; charset=utf-8")
	fmt.Fprintf(w, bio)
}

type Router struct{}

func (router Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	fmt.Println("path", path)
	switch path {
	case "/":
		homeHandler(w, r)
	case "/contact":
		contactHandler(w, r)
	case "/faq":
		faqHandler(w, r)
	case "/checkparam/":
		checkParamHandler(w, r)
	default:
		errNotFoundHandler(w, r)
	}
}

func main() {
	r := chi.NewRouter()
	// r.Use(middleware.Logger)
	r.Get("/", homeHandler)
	r.Route("/contact", func(api chi.Router) {
		api.Use(middleware.Logger)
		api.Get("/", contactHandler)
	})
	r.Get("/faq", faqHandler)
	r.Get("/checkparam/{param}", checkParamHandler)
	r.Get("/bio", bioHandler)
	r.NotFound(errNotFoundHandler)
	fmt.Println("Starting the server on :3000...")
	http.ListenAndServe(":3000", r)
}

func executeTemplate(w *http.ResponseWriter, temp_path templateDirAndPath, data interface{}) {
	tplPath := filepath.Join(temp_path.directory, temp_path.path)
	tpl, err := template.ParseFiles(tplPath)
	if err != nil {
		log.Printf("parsing template: %v", err)
		http.Error(*w,
			"There was an error parsing the template",
			http.StatusInternalServerError)
		return
	}

	err = tpl.Execute(*w, data)
	if err != nil {
		log.Printf("parsing template: %v", err)
		http.Error(*w,
			"There was an error parsing the template",
			http.StatusInternalServerError)
		return
	}
}
