package utils

import (
	"errors"
	"regexp"
)

func ValidateCEP(cep string) error {
	re := regexp.MustCompile(`^\d{8}$`)
	if !re.MatchString(cep) {
		return errors.New("invalid zipcode")
	}
	return nil
}
