package utils

import (
	"regexp"
	"strings"
)

func ToCode(rawValue string) (string, error) {
	rawValue = strings.TrimSpace(rawValue)
	reg, err := regexp.Compile("[^A-Za-z0-9]+")

	if err != nil {
		return "", err
	}

	code := reg.ReplaceAllString(strings.ToUpper(rawValue), "_")
	return code, nil
}

func IsValideMobileNumber(numberStr int) bool {

	return true

}
