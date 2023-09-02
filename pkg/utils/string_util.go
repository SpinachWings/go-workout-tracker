package utils

import "strings"

func StringIsBlank(s string) bool {
	return len(strings.Trim(s, " ")) == 0
}
