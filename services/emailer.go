package services

import (
	"embed"
	"fmt"
	"html/template"
	"io"
)

type Email struct {
	From        string
	To          string
	Content     string
	ContentType string
	Cc          []string
}

type Emailer interface {
	SendEmail(Email, io.Writer) error
}

type EmailTemplate struct {
	emailHTMLTpl *template.Template
}

func SendMail(mailer Emailer, email Email, writer io.Writer) error {
	err := mailer.SendEmail(email, writer)
	if err != nil {
		return err
	}
	return nil
}

var emailTplStrings = []string{
	"reset_password_email.gohtml",
}

//go:embed email_templates
var FS embed.FS

func loadEmailTemplates() (tpl *EmailTemplate) {
	tpl = &EmailTemplate{}
	loadedTemplate := template.New("base")
	templateStrings := getTemplatePaths(emailTplStrings, "email_templates")
	loadedTemplate = template.Must(loadedTemplate.ParseFS(FS, templateStrings...))
	tpl.emailHTMLTpl = loadedTemplate
	return tpl
}

func getTemplatePaths(tplStrings []string, baseFolderName string) []string {
	fullPaths := make([]string, len(tplStrings))
	for i, tplString := range tplStrings {
		fullPath := fmt.Sprintf("%s/%s", baseFolderName, tplString)
		fullPaths[i] = fullPath
	}
	return fullPaths
}
