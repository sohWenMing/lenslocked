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
		false,
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
		false,
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
		false,
	},
}

type ResetPasswordForm struct {
	NewPasswordInputAttribs     inputHTMLAttribs
	ConfirmPasswordInputAttribs inputHTMLAttribs
	ResetPasswordToken          string
}

var ResetPasswordFormData = ResetPasswordForm{
	NewPasswordInputAttribs: inputHTMLAttribs{
		"enter-password",
		"enter-password",
		"password",
		"Enter Password",
		"",
		"Enter Password",
		true,
		"",
		false,
	},
	ConfirmPasswordInputAttribs: inputHTMLAttribs{
		"confirm-password",
		"confirm-password",
		"password",
		"Confirm Password",
		"",
		"Confirm Password",
		true,
		"",
		false,
	},
	ResetPasswordToken: "",
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

type ResetPasswordTokenInfo struct {
	ResetPasswordToken string
}
