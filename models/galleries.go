package models

import (
	"database/sql"
	"fmt"
)

// Gallery houses fields that map to database structure that defines a gallery
type Gallery struct {
	ID     int
	UserID int
	Title  string
}

// Service that allows for gallery to have a connection to sql.DB methods, to be able to run database commands
type GalleryService struct {
	DB *sql.DB
}

// Creates a new gallery based on input title and userId. Returns pointer to a Gallery struct if successful, else
// returns nil and error
func (service *GalleryService) Create(title string, userId int) (*Gallery, error) {
	fmt.Println("user id entered: ", userId)

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
func (service *GalleryService) getByUserId(userId int) (*Gallery, error) {
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
