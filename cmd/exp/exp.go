package main

import (
	"html/template"
	"os"
)

type User struct {
	Name string
	Bio  string
	Age  int
	Gender
	TestValues
	StringToInt
}

type Gender struct {
	Name         string
	IsGay        bool
	IsHasPronoun bool
}

type TestValues struct {
	TestInts    []int
	TestStrings []string
}

type StringToInt map[string][]int

func main() {
	t, err := template.ParseFiles("hello.gohtml")
	if err != nil {
		panic(err)
	}

	user := User{
		Name: "Soh Wen Ming",
		Bio:  `<script>alert("Haha, you have been h4x0r3d!");</script>`,
		Age:  200,
		Gender: Gender{
			"male",
			true,
			false,
		},
		TestValues: TestValues{
			[]int{1, 2, 3, 4, 5},
			[]string{"one", "two", "three", "four", "five"},
		},
		StringToInt: StringToInt{
			"one to five": []int{1, 2, 3, 4, 5},
			"six to 10":   []int{6, 7, 8, 9, 10},
		},
	}

	err = t.Execute(os.Stdout, user)
	if err != nil {
		panic(err)
	}
}
