package internal

import (
	"errors"
	"os"
	"path"
)

func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func FileExists(s string) bool {
	if _, err := os.Stat(s); err == nil {
		return true
	} else if os.IsNotExist(err) {
		return false
	} else {
		return false
	}
}

func NeedSomeOfTheseFiles(folder string, acceptableFiles []string, notFoundErrorMessage string) (string, error) {
	for _, filename := range acceptableFiles {
		filepath := path.Join(folder, filename)
		if _, err := os.Stat(filepath); err == nil {
			return filepath, nil
		}
	}
	return "", errors.New(notFoundErrorMessage)
}
