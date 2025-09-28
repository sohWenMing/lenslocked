package models

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

// Gallery houses fields that map to database structure that defines a gallery
type Gallery struct {
	ID     int
	UserID int
	Title  string
}

// Service that allows for gallery to have a connection to sql.DB methods, to be able to run database commands
type GalleryService struct {
	DB        *sql.DB
	ImagesDir string
}

func (service *GalleryService) CreateImage(galleryId int, filename string, contents io.Reader) error {
	galleryDir := service.GalleryDir(galleryId)
	err := os.MkdirAll(galleryDir, 0755)
	if err != nil {
		return fmt.Errorf("creating gallery-%d images directory: %w", galleryId, err)
	}
	imagePath := filepath.Join(galleryDir, filename)
	dst, err := os.Create(imagePath)
	if err != nil {
		return fmt.Errorf("creating image file: %w", err)
	}
	defer dst.Close()
	_, err = io.Copy(dst, contents)
	if err != nil {
		return fmt.Errorf("copying contents to image: %w", err)
	}
	return nil
}

func ValidateContentType(r io.ReadSeeker, exts []string) error {
	bytesToValidate, err := get512Bytes(r)
	if err != nil {
		return err
	}
	contentType := http.DetectContentType(bytesToValidate)
	if !slices.Contains(exts, strings.ToLower(strings.TrimSpace(contentType))) {
		return errors.New("fileType not allowed")
	}
	return nil

}

func get512Bytes(r io.ReadSeeker) (bytesToValidate []byte, err error) {
	bytes := make([]byte, 512)
	_, err = r.Read(bytes)
	if err != nil {
		return []byte{}, nil
	}
	_, err = r.Seek(0, 0)
	if err != nil {
		return []byte{}, nil
	}
	return bytes, nil
}

// what i want is to get access to 512 bytes, for every file that I am trying ot read
// I also need to be able to reset the reading of the file, using seek
// Creates a new gallery based on input title and userId. Returns pointer to a Gallery struct if successful, else
// returns nil and error
func (service *GalleryService) GalleryDir(id int) string {
	imagesDir := service.ImagesDir
	if imagesDir == "" {
		imagesDir = "images"
	}
	return filepath.Join(imagesDir, fmt.Sprintf("%d", id))
}

func (service *GalleryService) GetAllowableContentTypes() []string {
	return []string{
		"image/jpeg",
		"image/jpg",
		"image/gif",
		"image/png",
	}
}

func (service *GalleryService) GetImageExtensions() []string {
	return []string{
		".png", ".gif", ".jpg", ".jpeg",
	}
}
func (service *GalleryService) Create(title string, userId int) (*Gallery, error) {

	gallery := Gallery{
		UserID: userId,
		Title:  title,
	}

	row := service.DB.QueryRow(
		`
		INSERT INTO galleries(title, user_id)
		VALUES ($1, $2)
		RETURNING id;
		`, title, userId,
	)

	err := row.Scan(
		&gallery.ID,
	)
	if err != nil {
		return nil, err
	}
	return &gallery, nil
}

// Deletes a gallery, based on the input galleryId. will return error if problem occurs else will return nil
func (service *GalleryService) DeleteById(galleryId int) error {
	_, err := service.DB.Exec(`
	DELETE from galleries
	WHERE id = ($1)	;
	`, galleryId)
	if err != nil {
		return err
	}
	return nil
}

func (service *GalleryService) GetById(galleryId int) (*Gallery, error) {
	row := service.DB.QueryRow(
		`SELECT galleries.id, galleries.user_id, galleries.title
		FROM galleries
		WHERE galleries.id = ($1)
		;
		`, galleryId,
	)
	var gallery Gallery
	err := row.Scan(&gallery.ID, &gallery.UserID, &gallery.Title)
	if err != nil {
		return nil, HandlePgError(err, &sqlNoRowsErrStruct{NoGalleryFound})
	}
	return &gallery, nil
}

func (service *GalleryService) GetGalleryListByUserId(userId int) ([]*Gallery, error) {
	rows, err := service.DB.Query(
		`SELECT galleries.id, galleries.user_id, galleries.title
		FROM galleries
		WHERE galleries.user_id = ($1)
		;
		`, userId,
	)
	returnedGalleries := []*Gallery{}
	if err != nil {
		return returnedGalleries, err
	}
	defer rows.Close()
	for rows.Next() {
		galleryToAppend := Gallery{}
		err := rows.Scan(&galleryToAppend.ID, &galleryToAppend.UserID, &galleryToAppend.Title)
		if err != nil {
			return []*Gallery{}, err
		}
		returnedGalleries = append(returnedGalleries, &galleryToAppend)
	}
	err = rows.Err()
	if err != nil {
		return []*Gallery{}, err
	}
	return returnedGalleries, err
}
func (service *GalleryService) GetByUserId(userId int) (*Gallery, error) {
	row := service.DB.QueryRow(
		`SELECT galleries.id, galleries.user_id, galleries.title
		FROM galleries
		WHERE galleries.user_id = ($1)
		;
		`, userId,
	)
	var gallery Gallery
	err := row.Scan(&gallery.ID, &gallery.UserID, &gallery.Title)
	if err != nil {
		return nil, HandlePgError(err, &sqlNoRowsErrStruct{NoGalleryFound})
	}
	return &gallery, nil
}
func (service *GalleryService) UpdateTitle(id int, title string) (err error) {
	result, err := service.DB.Exec(
		`
		UPDATE galleries
		SET title = ($1)
		WHERE id = ($2);
		`, title, id,
	)
	if err != nil {
		return fmt.Errorf("update gallery title %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("update gallery title %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no rows were affected - gallery id passed in: %d", id)
	}
	return nil
}
