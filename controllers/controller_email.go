package controllers

import "io"

type Email struct {
	From        string
	To          string
	Content     string
	ContentType string
	Cc          []string
}
type EmailService interface {
	SendEmail(Email, io.Writer) error
}
