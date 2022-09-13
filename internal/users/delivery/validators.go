package delivery

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"unicode"
	"user_service/internal/users"
)

const PASS_MIN_LENGTH = 5

func validateUser(user *users.UserIn) error {
	var errMsg []string

	if user.FirstName == "" {
		errMsg = append(errMsg, "firstName: may not be empty")
	}
	if user.LastName == "" {
		errMsg = append(errMsg, "lastName: may not be empty")
	}
	if user.Email == "" {
		errMsg = append(errMsg, "email: may not be empty")
	} else {
		if ok := isEmailValid(user.Email); !ok {
			errMsg = append(errMsg, "email: invalid")
		}
	}
	if user.Password == "" {
		errMsg = append(errMsg, "password: may not be empty")
	} else {
		verifier := NewPassVerifier(user.Password)
		verifier.Verify()
		if ok := verifier.IsValid(); !ok {
			for _, msg := range verifier.ErrorMessages() {
				errMsg = append(errMsg, fmt.Sprintf("password: %s", msg))
			}
		}
	}

	if len(errMsg) > 0 {
		return errors.New(strings.Join(errMsg, "; "))
	}
	return nil
}

func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(e)
}

type PassVerifier struct {
	// initial args
	password string
	// config
	minLength int
	// validation result
	length         bool
	containLetters bool
	containNumbers bool
	containUpper   bool
	containSpecial bool
}

//type Rule func(l rune) bool
//
//func MinLengthRule(min int) Rule {
//	cnt := 0
//	return func(l rune) bool {
//		cnt ++
//		return cnt >= min
//	}
//}
//
//func MinLengthRule(min int) Rule {
//	cnt := 0
//	return func(l rune) bool {
//		cnt ++
//		return cnt >= min
//	}
//}

func NewPassVerifier(password string) *PassVerifier {
	return &PassVerifier{password: password, minLength: PASS_MIN_LENGTH}
}

func (v *PassVerifier) IsValid() bool {
	return v.length && v.containLetters && v.containNumbers && v.containUpper && v.containSpecial
}

func (v *PassVerifier) Verify() {
	letters := 0
	for _, c := range v.password {
		switch {
		case unicode.IsNumber(c):
			v.containNumbers = true
		case unicode.IsUpper(c):
			v.containUpper = true
			letters++
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			v.containSpecial = true
		case unicode.IsLetter(c) || c == ' ':
			v.containLetters = true
		}
	}
	if len(v.password) > 5 {
		v.length = true
	}
}

func (v *PassVerifier) ErrorMessages() []string {
	var errMsg []string
	if !v.length {
		errMsg = append(errMsg, fmt.Sprintf("length should be greater then %d", v.minLength))
	}
	if !v.containLetters {
		errMsg = append(errMsg, fmt.Sprintf("should contain letters"))
	}
	if !v.containNumbers {
		errMsg = append(errMsg, fmt.Sprintf("should contain numbers"))
	}
	if !v.containUpper {
		errMsg = append(errMsg, fmt.Sprintf("should contain upper case letters"))
	}
	if !v.containSpecial {
		errMsg = append(errMsg, fmt.Sprintf("should contain special characters"))
	}
	return errMsg
}
