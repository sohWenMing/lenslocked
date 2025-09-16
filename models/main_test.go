package models

import (
	"fmt"
	"os"
	"testing"
)

var dbc *DBConnections
var baseUserEmailToPlainTextPassword = UserEmailToPlainTextPassword{"hello@test.com", "Holoq123holoq123"}

func TestMain(m *testing.M) {
	databaseConnection, err := InitDBConnections()
	if err != nil {
		fmt.Println("error occured during initialisation of db connection during test")
		os.Exit(1)
	}
	dbc = databaseConnection
	code := m.Run()
	if err := dbc.DB.Close(); err != nil {
		fmt.Println("error closing db: ", err)
	}
	os.Exit(code)
}
