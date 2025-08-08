package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		cookie := &http.Cookie{
			Name:  "TestCookie",
			Value: "This-is-a-test-cookie",
			Path:  "/test-cookie",
		}
		http.SetCookie(w, cookie)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("check your browser for a cookie!"))
	})
	fmt.Println("server listening on port 4000")
	log.Fatal(http.ListenAndServe(":4000", r))
}

// i want a function that adds a prefix to a string
