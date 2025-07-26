package controllers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestGetUrlParam(t *testing.T) {
	type test struct {
		name           string
		fullUrlRequest string
		baseUrl        string
		want           string
		isExpectErr    bool
	}

	tests := []test{
		{
			"test get persona param",
			"/test-request/persona",
			"/test-request",
			"persona",
			false,
		},
	}
	for i, test := range tests {
		req, err := http.NewRequest(http.MethodGet, test.fullUrlRequest, nil)
		if err != nil {
			t.Errorf("test failed on test %d: didn't expect error, got %v", err, i)
			return
		}

		router := chi.NewRouter()
		paramFunc := func(param string) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				got, err := GetUrlParam(r, param)
				switch test.isExpectErr {
				case true:
					if err == nil {
						t.Errorf("test failed on test %d: expected error, didn't get one", i)
					}
				case false:
					if err != nil {
						t.Errorf("test failed on test %d: didn't expect error, got %v", i, err)
					}
					if got != test.want {
						t.Errorf("got %s, want %s", got, test.want)
					}
				}
			}
		}
		fullTestUrl := fmt.Sprintf("%s/{%s}", test.baseUrl, test.want)
		router.Get(fullTestUrl, paramFunc(test.want))
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
	}

}
