package integrationtesting

import (
	"bytes"
	"io"
	"mime/quotedprintable"
	"net/mail"
	"testing"

	"github.com/sohWenMing/lenslocked/services"
)

func TestEmailGeneration(t *testing.T) {
	fromEmail := "wenming.soh@gmail.com"
	toEmail := "sarahlinshuyi@gmail.com"
	content := `This is a text email with a <a href="http://www.google.com">link</a>`
	buf := bytes.Buffer{}
	err := mailer.SendEmail(services.Email{
		From:        fromEmail,
		To:          toEmail,
		Content:     content,
		ContentType: "text/html",
		Cc:          []string{},
	}, &buf)
	if err != nil {
		t.Errorf("didn't expect err, got %v\n", err)
	}

	// what is being written to the buffer is a RFC 5322 / MIME string, so we can parse it with net/mail ReadMessage function
	readMessage, err := mail.ReadMessage(&buf)
	if err != nil {
		t.Errorf("didn't expect err, got %v\n", err)
	}
	// the struct returned from ReadMessage function has a Header, which is a map that has some methods attached, and the body
	// which is an io.Reader
	from := readMessage.Header.Get("From")
	to := readMessage.Header.Get("To")
	//we want to wrap the body in a quotedPrintable wrapper, because we want to be able to decode the quoted printable encoding
	quotedPrintableWrapper := quotedprintable.NewReader(readMessage.Body)
	// read all the bytes from the reader, which returns slice of bytes and error (if there is an error)
	readBody, err := io.ReadAll(quotedPrintableWrapper)
	if err != nil {
		t.Errorf("didn't expect err, got %v\n", err)
	}
	readBodyString := string(readBody)
	if from != fromEmail {
		t.Errorf("got %s, want %s\n", from, fromEmail)
	}
	if to != toEmail {
		t.Errorf("got %s, want %s\n", to, toEmail)
	}
	if readBodyString != content {
		t.Errorf("got %s, want %s\n", readBodyString, content)
	}
}
