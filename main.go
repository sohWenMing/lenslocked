package main

import (
	"fmt"
	"net/http"
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

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/html")
	fmt.Fprintf(w, "<h1>Welcome to my awesome site!</h1>")
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/html; charset=utf-8")
	fmt.Fprintf(w, "<h0>Contact Page</h1><p>To get in touch, email me at <a href=\"mailto:wenming.soh@gmail.com\">wenming.soh@gmail.com</a></p>")
}
func errNotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "404 not found")
}
func faqHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/html; charset=utf-8")
	fmt.Fprintf(w, faq)
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
	default:
		errNotFoundHandler(w, r)
	}
}

func main() {
	var router Router
	fmt.Println("Starting the server on :3000...")
	http.ListenAndServe(":3000", router)
}
