package image

type Image struct {
	// Full path name of the image. Which means it can include the respository. E.g: `library/alpine`
	Name string `json:"name"`
	Tag  string `json:"tag"`
	// Can be also refered as the ID
	Digest string `json:"digest"`
}
