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

type EmailData struct {
	URL string
}

type EmailService struct {
	Emailer
	*EmailTemplate
}

func (e *EmailService) SendMail(email Email, writer io.Writer) error {
	err := e.SendEmail(email, writer)
	if err != nil {
		return err
	}
	return nil
}

func InitEmailService(emailer Emailer, emailTemplate *EmailTemplate) *EmailService {
	return &EmailService{
		emailer, emailTemplate,
	}
}

type Emailer interface {
	SendEmail(Email, io.Writer) error
}

type EmailTemplate struct {
	EmailHTMLTpl *template.Template
}

var emailTplStrings = []string{
	"reset_password_email.gohtml",
}

//go:embed email_templates
var FS embed.FS

func LoadEmailTemplates() (tpl *EmailTemplate) {
	tpl = &EmailTemplate{}
	loadedTemplate := template.New("base")
	templateStrings := getTemplatePaths(emailTplStrings, "email_templates")
	loadedTemplate = template.Must(loadedTemplate.ParseFS(FS, templateStrings...))
	tpl.EmailHTMLTpl = loadedTemplate
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
