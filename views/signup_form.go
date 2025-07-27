package views

type inputHTMLAttribs struct {
	Name         string
	Id           string
	InputType    string
	PlaceHolder  string
	AutoComplete string
	LabelText    string
	IsRequired   bool
}

type SignUpFormAttibs struct {
	EmailInputAttribs, PasswordInputAttribs inputHTMLAttribs
}

var SignUpFormData = SignUpFormAttibs{
	EmailInputAttribs: inputHTMLAttribs{
		"email",
		"email",
		"email",
		"Email Address",
		"email",
		"Email",
		true,
	},
	PasswordInputAttribs: inputHTMLAttribs{
		"password",
		"password",
		"password",
		"Password",
		"",
		"Password",
		true,
	},
}
