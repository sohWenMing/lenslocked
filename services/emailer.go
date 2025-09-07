package services

import "io"

// so what i need to do is to write a an interface, which allows injection of html

type Email struct {
	From        string
	To          string
	Content     string
	ContentType string
	Cc          []string
}

type EmailContentGenerator interface {
	GenerateEmailContent() (string, error)
}
type Emailer interface {
	SendEmail(Email, io.Writer) error
}

func SendMail(mailer Emailer, email Email, writer io.Writer) error {
	err := mailer.SendEmail(email, writer)
	if err != nil {
		return err
	}
	return nil
}
