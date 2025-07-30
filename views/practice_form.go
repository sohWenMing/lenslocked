package views

type PracticeForm struct {
	Input1Attribs, Input2Attribs inputHTMLAttribs
	CheckBoxAttribs              checkBoxHTMLAttribs
	ButtonAttribs                buttonAttribs
}

var PracticeFormData = PracticeForm{
	inputHTMLAttribs{
		Name:        "first_name",
		Id:          "first_name_input",
		InputType:   "text",
		PlaceHolder: "Enter your first name here",
		LabelText:   "First Name",
		IsRequired:  true,
	},
	inputHTMLAttribs{
		Name:        "last_name",
		Id:          "last_name_input",
		InputType:   "text",
		PlaceHolder: "Enter your last name here",
		LabelText:   "Last Name",
		IsRequired:  true,
	},
	checkBoxHTMLAttribs{
		"testCheckBox",
		"testCheckBox",
		"Are You a Human???",
		false,
	},
	buttonAttribs{
		"Submit",
	},
}
