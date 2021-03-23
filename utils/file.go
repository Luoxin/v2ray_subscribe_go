package utils

import (
	"io/ioutil"
	"os"
)

func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func IsFile(path string) bool {
	return !IsDir(path)
}

func CopyFile(from, to string) error {
	input, err := ioutil.ReadFile(from)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(to, input, 0644)
	if err != nil {
		return err
	}

	return nil
}
