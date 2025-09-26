package models

import (
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"
	"slices"
)

func LoadImageFileServer(path string) http.Handler {
	return http.FileServer(http.Dir(path))
}

type GalleryImage struct {
	galleryId       int
	path            string
	fileNameEscaped string
}

func (g *GalleryImage) GetPath() string {
	return g.path
}

func GetImagesByGalleryId(galleryId int, exts []string) (galleryImages []*GalleryImage, err error) {
	globPattern := fmt.Sprintf("./images/%d/*", galleryId)
	filepaths, err := getImagePaths(globPattern, exts)
	if err != nil {
		return []*GalleryImage{}, err
	}
	returnedGalleryImages := make([]*GalleryImage, len(filepaths))
	for i, filePath := range filepaths {
		fileNameEscaped := url.PathEscape(filepath.Base(filePath))
		returnedGalleryImages[i] = &GalleryImage{
			galleryId:       galleryId,
			path:            (fmt.Sprintf("/galleries/%d/images/%s", galleryId, fileNameEscaped)),
			fileNameEscaped: fileNameEscaped,
		}
	}
	return returnedGalleryImages, nil
}

func getImagePaths(globPattern string, exts []string) (filepaths []string, err error) {
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
