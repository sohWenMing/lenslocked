package views

type SignUpForm struct {
	EmailInputAttribs, PasswordInputAttribs inputHTMLAttribs
}

var SignUpFormData = SignUpForm{
	EmailInputAttribs: inputHTMLAttribs{
		"email",
		"email",
		"email",
		"Email Address",
		"email",
		"Email",
		true,
		"",
		// Value can be set by a method
	},
	PasswordInputAttribs: inputHTMLAttribs{
		"password",
		"password",
		"password",
		"Password",
		"",
		"Password",
		true,
		"",
	},
}

func (s *SignUpForm) SetEmailValue(input string) {
	s.EmailInputAttribs.Value = input
}
func (s *SignUpForm) SetPasswordValue(input string) {
	s.PasswordInputAttribs.Value = input
}
