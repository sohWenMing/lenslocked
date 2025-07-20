package models

type User struct {
	Name   string
	Age    int
	Skills []string
}

var UserMap = map[string]User{
	"wen": {
		"Wen", 42, []string{"management", "coding", "negotiation"},
	},
	"sarah": {
		"Sarah", 38, []string{"pilates", "mothering", "gyrotonic"},
	},
}
