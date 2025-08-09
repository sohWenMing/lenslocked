package models

import "testing"

func TestHashSessionToken(t *testing.T) {
	type test struct {
		testname       string
		token          string
		isExpectMatch  bool
		isExpectErr    bool
		isTestNullHash bool
	}

	tests := []test{
		{
			"basic test - simple token",
			"token1234",
			true,
			false,
			false,
		},
		{
			"fail test - null token",
			"",
			false,
			true,
			false,
		},
		{
			"fail test - null hash",
			"token1234",
			false,
			true,
			true,
		},
	}

	for _, test := range tests {
		t.Run(test.testname, func(t *testing.T) {
			initialHash := HashSessionToken(test.token)
			if test.isTestNullHash {
				initialHash = ""
			}
			isVerified, err := VerifySessionToken(test.token,
				initialHash)
			switch test.isExpectErr {
			case true:
				if err == nil {
					t.Errorf("expected error, didn't get one")
				}
			default:
				if err != nil {
					t.Errorf("didn't expect err, got %v\n", err)
				}
			}
			if isVerified != test.isExpectMatch {
				t.Errorf("got %t want %t", isVerified, test.isExpectMatch)
			}
		})
	}
}

func TestVerifySessionToken(t *testing.T) {
	token, err := SessionToken()
	if err != nil {
		t.Errorf("didn't expect err, got %v\n", err)
	}
	initialHash := HashSessionToken(token)
	isVerified, err := VerifySessionToken(token, initialHash)
	if !isVerified {
		t.Errorf("got %t, want %t", isVerified, true)
	}

}
