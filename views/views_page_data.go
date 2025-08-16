package views

import "html/template"

type SignInSignUpForm struct {
	EmailInputAttribs, PasswordInputAttribs inputHTMLAttribs
	CSRFField                               template.HTML
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
func (s *SignInSignUpForm) SetCSRFFormValue(input template.HTML) {
	s.CSRFField = input
}
