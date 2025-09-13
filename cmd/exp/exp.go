package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
)

func main() {
	handler := func(w http.ResponseWriter, r *http.Request) {
		requestURL := r.URL
		fmt.Println("url: ", requestURL)
		path := requestURL.Path
		fmt.Println("path:", path)
		rawQuery := requestURL.RawQuery
		fmt.Println("rawQuery: ", rawQuery)

		values, err := url.ParseQuery(rawQuery)
		if err != nil {
			fmt.Println("didn't expect error, got %v", err)
		}
		fmt.Println("values", values)

	}
	url := "http://localhost:3000/reset_password?token=fbdcc1c0-2903-49a8-8f98-e56c9be9e353"
	request := httptest.NewRequest(http.MethodGet, url, nil)
	w := httptest.NewRecorder()
	handler(w, request)

}
