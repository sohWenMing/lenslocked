package models

import "net/http"

func LoadImageFileServer(path string) http.Handler {
	return http.FileServer(http.Dir(path))
}
