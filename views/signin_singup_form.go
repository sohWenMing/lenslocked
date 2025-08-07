package views

type SignInSignUpForm struct {
	EmailInputAttribs, PasswordInputAttribs inputHTMLAttribs
}

var SignUpSignInFormData = SignInSignUpForm{
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

func (s *SignInSignUpForm) SetEmailValue(input string) {
	s.EmailInputAttribs.Value = input
}
func (s *SignInSignUpForm) SetPasswordValue(input string) {
	s.PasswordInputAttribs.Value = input
}
