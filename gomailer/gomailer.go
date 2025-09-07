package gomailer

import (
	"io"

	"github.com/sohWenMing/lenslocked/services"
	"gopkg.in/gomail.v2"
)

type GoMailer struct {
	dialer *gomail.Dialer
}

func NewGoMailer(host, username, password string, port int) *GoMailer {
	return &GoMailer{
		gomail.NewDialer(host, port, username, password),
	}
}

func (g *GoMailer) SendEmail(email services.Email, w io.Writer) error {
	m := g.PrepEmail(email, w)
	if w != nil {
		m.WriteTo(w)
	}
	err := g.dialer.DialAndSend(m)
	if err != nil {
		return err
	}
	return nil
}

func (g *GoMailer) PrepEmail(email services.Email, w io.Writer) *gomail.Message {
	m := gomail.NewMessage()
	m.SetHeader("From", email.From)
	m.SetHeader("To", email.To)
	if len(email.Cc) > 0 {
		m.SetHeader("Cc", email.Cc...)
	}
	m.SetBody(email.ContentType, email.Content)
	return m
}
