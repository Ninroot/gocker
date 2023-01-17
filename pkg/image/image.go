package image

import (
	"errors"
	"strings"
)

type Image struct {
	// Full path name of the image. Which means it can include the respository. E.g: `library/alpine`
	Name string `json:"name"`
	Tag  string `json:"tag"`
	// Can be also refered as the ID
	Digest string `json:"digest"`
}

// name[:TAG]
// return repository and tag
func Parse(imageName string) (Image, error) {
	s := strings.Split(imageName, ":")
	if len(s) == 1 {
		return Image{Name: s[0], Tag: "latest"}, nil
	}
	if len(s) == 2 {
		return Image{Name: s[0], Tag: s[1]}, nil
	}
	return Image{}, errors.New("image name has the wrong format")
}
