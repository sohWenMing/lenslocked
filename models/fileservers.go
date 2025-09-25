package models

import (
	"net/http"
	"path/filepath"
	"slices"
)

func LoadImageFileServer(path string) http.Handler {
	return http.FileServer(http.Dir(path))
}

func GetImagePaths(globPattern string, exts []string) (filepaths []string, err error) {
	returnedPaths := []string{}
	files, err := filepath.Glob(globPattern)
	if err != nil {
		return returnedPaths, err
	}

	for _, file := range files {
		if slices.Contains(exts, filepath.Ext(file)) {
			returnedPaths = append(returnedPaths, file)
		}
	}
	return returnedPaths, nil
}
