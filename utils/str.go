package utils

import (
	"strings"
)

func ShortStr4Web(str string, max int) string {
	str = strings.ReplaceAll(str, "\n", "\\n")
	str = strings.ReplaceAll(str, "\r", "\\r")
	str = strings.ReplaceAll(str, "\t", "\\t")
	return ShortStr(str, max)
}

func ShortStr(str string, max int) string {
	if len(str) > max {
		return str[:max] + "..."
	}
	return str
}
