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

type ForgotPasswordForm struct {
	EmailInputAttribs inputHTMLAttribs
}

var ForgotPasswordFormData = ForgotPasswordForm{
	EmailInputAttribs: inputHTMLAttribs{
		"email",
		"email",
		"email",
		"Email Address",
		"email",
		"Email",
		true,
		"",
	},
}

type PageData struct {
	UserId    int
	OtherData any
}

func InitPageData(userId int, otherData any) PageData {
	return PageData{
		UserId:    userId,
		OtherData: otherData,
	}
}

func (s *SignInSignUpForm) SetEmailValue(input string) {
	s.EmailInputAttribs.Value = input
}
func (s *SignInSignUpForm) SetPasswordValue(input string) {
	s.PasswordInputAttribs.Value = input
}
