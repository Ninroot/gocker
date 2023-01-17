package util

import (
	"log"
	"os"
)

func EnsureDir(path string) string {
	if err := os.MkdirAll(path, 0775); err != nil {
		log.Println("Could not ensure directory:", path)
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
