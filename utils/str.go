package utils

import (
	"strings"
)

func ShortStr4Web(str string, max int) string {
	str = strings.ReplaceAll(str, "\n", "\\n")
	str = strings.ReplaceAll(str, "\r", "\\r")
	str = strings.ReplaceAll(str, "\t", "\\t")

	newStr := ShortStr(str, max)
	if len(newStr) != len(str) {
		newStr += "..."
	}
	return newStr
}

func ShortStr(str string, max int) string {
	if max < 0 {
		return str
	}
	if len(str) > max {
		return str[:max]
	}
	return str
}
