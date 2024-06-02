package utils

import "regexp"

func ValidatePassword(password string) bool {
	var (
		uppercase   = regexp.MustCompile(`[A-Z]`)
		lowercase   = regexp.MustCompile(`[a-z]`)
		number      = regexp.MustCompile(`[0-9]`)
		specialChar = regexp.MustCompile(`[!@#~$%^&*(),.?":{}|<>]`)
	)

	if len(password) < 8 {
		return false
	}

	if !uppercase.MatchString(password) || !lowercase.MatchString(password) || !number.MatchString(password) || !specialChar.MatchString(password) {
		return false
	}

	return true
}
