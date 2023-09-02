package utils

import "strings"

func IsMismatchedHashAndPassword(err error) bool {
	if err == nil {
		return false
	}
	return strings.HasSuffix(err.Error(), "hashedPassword is not the hash of the given password")
}
