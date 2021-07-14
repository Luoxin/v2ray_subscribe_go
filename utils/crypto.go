package utils

import (
	"crypto/md5"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
)

func Md5(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

func Sha256(s string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(s)))
}

func Sha224(s string) string {
	return fmt.Sprintf("%x", sha256.Sum224([]byte(s)))
}

func Sha512(s string) string {
	return fmt.Sprintf("%x", sha512.Sum512([]byte(s)))
}

func Sha384(s string) string {
	return fmt.Sprintf("%x", sha512.Sum384([]byte(s)))
}
