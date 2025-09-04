package gomailer

import (
	"io"

	"github.com/sohWenMing/lenslocked/services"
	"gopkg.in/gomail.v2"
)

type GoMailer struct {
	dialer *gomail.Dialer
	writer io.Writer
}

func NewGoMailer(host, username, password string, port int, writer io.Writer) *GoMailer {
	return &GoMailer{
		gomail.NewDialer(host, port, username, password),
		writer,
	}
}

func (g *GoMailer) SendEmail(email services.Email) error {
	m := gomail.NewMessage()
	m.SetHeader("From: ", email.From)
	m.SetHeader("To: ", email.To)
	if len(email.Cc) > 0 {
		m.SetHeader("Cc", email.Cc...)
	}
	m.SetBody(email.ContentType, email.ContentType)
	if g.writer != nil {
		m.WriteTo(g.writer)
	}
	err := g.dialer.DialAndSend(m)
	if err != nil {
		return err
	}
	return nil

}
