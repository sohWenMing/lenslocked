package models

import (
	"fmt"
	"testing"
)

func TestGalleryCreate(t *testing.T) {
	userIdToSession, shouldReturn := CreateTestUser(t)
	if shouldReturn {
		return
	}
	defer dbc.UserService.DeleteUserAndSession(userIdToSession.UserID)

	type test struct {
		name          string
		userId        int
		galleryName   string
		isErrExpected bool
	}

	tests := []test{
		{
			"passing gallery create test",
			userIdToSession.UserID,
			"passing gallery",
			false,
		},
		{
			"failing gallery create test",
			0,
			"failing gallery",
			true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gallery, err := dbc.GalleryService.Create(test.galleryName, test.userId)
			switch test.isErrExpected {
			case true:
				if err == nil {
					t.Errorf("expected error, didn't get one")
				}

			case false:
				if err != nil {
					t.Errorf("didn't expect error, got %v\n", err)
					return
				}

				if gallery.Title != test.galleryName {
					t.Errorf("got %s, want %s\n", gallery.Title, test.galleryName)
				}
				if gallery.UserID != test.userId {
					t.Errorf("got %s, want %s\n", gallery.Title, test.galleryName)
				}

				err := dbc.GalleryService.DeleteById(gallery.ID)
				if err != nil {
					t.Errorf("didn't expect error, got %v\n", err)
					return
				}
			}
		})
	}
}
func TestGetGalleryById(t *testing.T) {
	userIdToSession, shouldReturn := CreateTestUser(t)
	if shouldReturn {
		return
	}
	fmt.Println("userId: ", userIdToSession.UserID)
	defer dbc.UserService.DeleteUserAndSession(userIdToSession.UserID)
	//we're making the assumption here that the create function should work without problems, so we just create one to run the retrieval tests
	gallery, err := dbc.GalleryService.Create("test_gallery", userIdToSession.UserID)
	if err != nil {
		t.Errorf("didn't expect error, got %v\n", err)
		return
	}
	defer dbc.GalleryService.DeleteById(gallery.ID)

	type test struct {
		name                  string
		expectedErrMsg        string
		expectedGalleryTitle  string
		expectedGalleryId     int
		expectedGalleryUserId int
		galleryId             int
		isErrExpected         bool
	}

	tests := []test{
		{
			name:                  "test should pass",
			expectedErrMsg:        "",
			expectedGalleryId:     gallery.ID,
			expectedGalleryTitle:  "test_gallery",
			expectedGalleryUserId: userIdToSession.UserID,
			galleryId:             gallery.ID,
			isErrExpected:         false,
		},
		{
			name:                  "test should fail because of 0 galleryId",
			expectedErrMsg:        "no gallery was found with the input gallery Id",
			expectedGalleryId:     0,
			expectedGalleryTitle:  "test_gallery",
			expectedGalleryUserId: userIdToSession.UserID,
			galleryId:             0,
			isErrExpected:         true,
		},
		{
			name:                  "test should fail because of non-existent galleryId",
			expectedErrMsg:        "no gallery was found with the input gallery Id",
			expectedGalleryId:     0,
			expectedGalleryTitle:  "test_gallery",
			expectedGalleryUserId: userIdToSession.UserID,
			galleryId:             gallery.ID + 1,
			isErrExpected:         true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gallery, err := dbc.GalleryService.GetById(test.galleryId)
			switch test.isErrExpected {
			case true:
				if err == nil {
					t.Errorf("expected error, didn't get one")
				}
			case false:
				if err != nil {
					fmt.Println("error occured on GetById")
					t.Errorf("didn't expect error, got %v\n", err)
					return
				}
				if gallery.ID != test.expectedGalleryId {
					t.Errorf("got %d, want %d", gallery.ID, test.expectedGalleryId)
				}
				if gallery.UserID != test.expectedGalleryUserId {
					t.Errorf("got %d, want %d", gallery.ID, test.expectedGalleryId)
				}
				if gallery.Title != test.expectedGalleryTitle {
					t.Errorf("got %d, want %d", gallery.ID, test.expectedGalleryId)
				}
			}
		})
	}

}
func TestGetGalleryByUserId(t *testing.T) {
	userIdToSession, shouldReturn := CreateTestUser(t)
	if shouldReturn {
		return
	}
	defer dbc.UserService.DeleteUserAndSession(userIdToSession.UserID)
	//we're making the assumption here that the create function should work without problems, so we just create one to run the retrieval tests
	gallery, err := dbc.GalleryService.Create("test_gallery", userIdToSession.UserID)
	if err != nil {
		t.Errorf("didn't expect error, got %v\n", err)
		return
	}
	defer dbc.GalleryService.DeleteById(gallery.ID)

	type test struct {
		name                  string
		expectedErrMsg        string
		expectedGalleryTitle  string
		expectedGalleryId     int
		expectedGalleryUserId int
		galleryUserId         int
		isErrExpected         bool
	}

	tests := []test{
		{
			name:                  "test should pass",
			expectedErrMsg:        "",
			expectedGalleryId:     gallery.ID,
			expectedGalleryTitle:  "test_gallery",
			expectedGalleryUserId: userIdToSession.UserID,
			galleryUserId:         gallery.UserID,
			isErrExpected:         false,
		},
		{
			name:                  "test should fail because of 0 gallery user id",
			expectedErrMsg:        "no gallery was found with the input gallery Id",
			expectedGalleryId:     0,
			expectedGalleryTitle:  "test_gallery",
			expectedGalleryUserId: userIdToSession.UserID,
			galleryUserId:         0,
			isErrExpected:         true,
		},
		{
			name:                  "test should fail because of non-existent gallery user Id",
			expectedErrMsg:        "no gallery was found with the input gallery Id",
			expectedGalleryId:     0,
			expectedGalleryTitle:  "test_gallery",
			expectedGalleryUserId: userIdToSession.UserID,
			galleryUserId:         gallery.UserID + 1,
			isErrExpected:         true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gallery, err := dbc.GalleryService.getByUserId(test.galleryUserId)
			switch test.isErrExpected {
			case true:
				if err == nil {
					t.Errorf("expected error, didn't get one")
				}
			case false:
				if err != nil {
					fmt.Println("error occured on getByUserId")
					t.Errorf("didn't expect error, got %v\n", err)
					return
				}
				if gallery.ID != test.expectedGalleryId {
					t.Errorf("got %d, want %d", gallery.ID, test.expectedGalleryId)
				}
				if gallery.UserID != test.expectedGalleryUserId {
					t.Errorf("got %d, want %d", gallery.ID, test.expectedGalleryId)
				}
				if gallery.Title != test.expectedGalleryTitle {
					t.Errorf("got %d, want %d", gallery.ID, test.expectedGalleryId)
				}
			}
		})
	}

}

func CreateTestUser(t *testing.T) (*UserIdToSession, bool) {
	testUserEmailToPlainTextPassword := UserEmailToPlainTextPassword{
		"test_user@gmail.com",
		"Holoq123holoq123",
	}
	userIdToSession, err := dbc.UserService.CreateUser(testUserEmailToPlainTextPassword)
	if err != nil {
		t.Errorf("didn't expect error, got %v", err)
		return nil, true
	}
	return userIdToSession, false
}
