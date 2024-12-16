package email

import (
	"MEDODS-test/internal/domain/model"
	"errors"
	"regexp"
)

const (
	minEmailLen = 3
	maxEmailLen = 255
)

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func IsEmailValid(email string) bool {
	if len(email) < minEmailLen || len(email) > maxEmailLen {
		return false
	}

	return emailRegex.MatchString(email)
}

func ValidateEmailInput(input model.SendEmailInput) error {
	if input.To == "" {
		return errors.New("empty to")
	}

	if input.Subject == "" || input.Body == "" {
		return errors.New("empty subject/body")
	}

	if !IsEmailValid(input.To) {
		return errors.New("invalid to email")
	}

	return nil
}
