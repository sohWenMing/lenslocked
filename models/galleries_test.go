package models

import "testing"

func TestGalleryCreate(t *testing.T) {
	testUserEmailToPlainTextPassword := UserEmailToPlainTextPassword{
		"test_user@gmail.com",
		"Holoq123holoq123",
	}
	userIdToSession, err := dbc.UserService.CreateUser(testUserEmailToPlainTextPassword)
	if err != nil {
		t.Errorf("didn't expect error, got %v", err)
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
				err := dbc.GalleryService.DeleteById(gallery.ID)
				if err != nil {
					t.Errorf("didn't expect error, got %v\n", err)
					return
				}
			}
		})
	}
}
