package pkg

import (
	"errors"
	"log"
	"os"
	"strings"
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

// name[:TAG]
// return repository and tag
func Parse(imageName string) (ImageId, error) {
	s := strings.Split(imageName, ":")
	if len(s) == 1 {
		return ImageId{Name: s[0], Tag: "latest"}, nil
	}
	if len(s) == 2 {
		return ImageId{Name: s[0], Tag: s[1]}, nil
	}
	return ImageId{}, errors.New("image name has the wrong format")
}
