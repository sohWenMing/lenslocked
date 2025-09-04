package services

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
	SendEmail(Email) error
}

func SendEmail(mailer Emailer, email Email) error {
	err := mailer.SendEmail(email)
	if err != nil {
		return err
	}
	return nil
}
