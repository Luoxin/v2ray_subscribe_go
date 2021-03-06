package utils

import (
	"encoding/base64"
)

func Base64DecodeStripped(s string) (string, error) {
	// if i := len(s) % 4; i != 0 {
	// 	s += strings.Repeat("=", 4-i)
	// }
	// s = strings.ReplaceAll(s, " ", "+")
	decoded, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		decoded, err = base64.URLEncoding.DecodeString(s)
		if err != nil {
			decoded, err = base64.RawStdEncoding.DecodeString(s)
			if err != nil {
				decoded, err = base64.RawURLEncoding.DecodeString(s)
			}
		}
	}
	return string(decoded), err
}
