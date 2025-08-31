package main

import (
	"os"

	"gopkg.in/gomail.v2"
)

func main() {
	m := gomail.NewMessage()
	m.SetHeader("From", "wenming.soh@nindgabeet.com")
	m.SetHeader("To", "wenming.soh@gmail.com")
	m.SetAddressHeader("Cc", "sarahlinshuyi@gmail.com", "Sarah")
	m.SetBody("text/html", "Hello <b>Wen</b")
	m.WriteTo(os.Stdout)

	d := gomail.NewDialer("sandbox.smtp.mailtrap.io", 587, "ec07c285658e45", "b633f2509083f8")
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}

}
