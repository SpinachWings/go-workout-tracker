package utils

import "strings"

func StandardiseCase(s string) string {
	trimmedString := strings.Trim(s, " ")
	firstChar := strings.ToUpper(trimmedString[:1])
	restOfString := strings.ToLower(trimmedString[1:])
	return firstChar + restOfString
}
