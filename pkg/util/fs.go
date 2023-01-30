package util

import (
	"os"

	"github.com/sirupsen/logrus"
)

func EnsureDir(path string) string {
	if err := os.MkdirAll(path, 0775); err != nil {
		logrus.WithField("path", path).Error("Could not ensure directory")
	}
	return path
}

// returns wether the item exist in the path
func Exist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
