package services

import (
	"bytes"
	"fmt"
	"testing"
)

func TestResetPasswordTemplate(t *testing.T) {
	type emailData struct {
		URL string
	}
	testData := emailData{
		"https://www.google.com",
	}
	buf := bytes.Buffer{}
	emailTemplate := loadEmailTemplates()
	err := emailTemplate.emailHTMLTpl.ExecuteTemplate(&buf, "reset_password_email.gohtml", testData)
	if err != nil {
		t.Errorf("didn't expect error, got %v\n", err)
	}
	fmt.Println("returned string:", buf.String())
}
