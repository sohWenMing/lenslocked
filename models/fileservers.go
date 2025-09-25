package models

import (
	"fmt"
	"net/http"
	"path/filepath"
	"slices"
)

func LoadImageFileServer(path string) http.Handler {
	return http.FileServer(http.Dir(path))
}

type GalleryImage struct {
	galleryId int
	path      string
	fileName  string
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
		fileName := filepath.Base(filePath)
		returnedGalleryImages[i] = &GalleryImage{
			galleryId: galleryId,
			path:      fmt.Sprintf("/galleries/%d/images/%s", galleryId, fileName),
			fileName:  fileName,
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
