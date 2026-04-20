package utils

import "net/mail"

// IsValidEmail checks whether the input has a valid email address format.
func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
