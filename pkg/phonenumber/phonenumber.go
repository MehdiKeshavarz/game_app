package phonenumber

import "regexp"

func IsValid(phoneNumber string) bool {
	pattern := regexp.MustCompile(`^(0|0098|\+98)9(0[1-5]|[1 3]\d|2[0-2]|98)\d{7}$`)

	if !pattern.MatchString(phoneNumber) {
		return false
	}

	return true
}
