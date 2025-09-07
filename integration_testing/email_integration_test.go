package integrationtesting

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/sohWenMing/lenslocked/services"
)

func TestEmailGeneration(t *testing.T) {
	buf := bytes.Buffer{}
	err := mailer.SendEmail(services.Email{
		From:        "wenming.soh@gmail.com",
		To:          "sarahlinshuyi@gmail.com",
		Content:     `This is a text email with a <a href="http://www.google.com">link</a>`,
		ContentType: "text/html",
		Cc:          []string{},
	}, &buf)
	if err != nil {
		t.Errorf("didn't expect err, got %v\n", err)
	}
	fmt.Println(buf.String())
}
