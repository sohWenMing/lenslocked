package integrationtesting

import (
	"log"
	"os"
	"testing"

	"github.com/sohWenMing/lenslocked/models"
)

var isDev = false

func TestMain(m *testing.M) {
	envVars, err := models.LoadEnv("../.env")
	if err != nil {
		log.Fatal(err)
	}
	isDev = envVars.IsDev
	code := m.Run()
	os.Exit(code)
}

func TestEnvLoading(t *testing.T) {
	want := true
	got := isDev
	if got != want {
		t.Errorf("Got %v, want %v\n", got, want)
	}
}
