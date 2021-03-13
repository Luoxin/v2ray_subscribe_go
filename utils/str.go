package utils

func ShortStr(str string, max int) string {
	if len(str) > max {
		return str[:max]
	}
	return str
}
