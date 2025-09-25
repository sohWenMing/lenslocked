package models

import (
	"database/sql"
	"fmt"
	"path/filepath"
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

// Creates a new gallery based on input title and userId. Returns pointer to a Gallery struct if successful, else
// returns nil and error
func (service *GalleryService) GalleryDir(id int) string {
	imagesDir := service.ImagesDir
	if imagesDir == "" {
		imagesDir = "images"
	}
	return filepath.Join(imagesDir, fmt.Sprintf("gallery-%d", id))
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
