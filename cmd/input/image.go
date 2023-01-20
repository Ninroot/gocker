package input

import (
	"strings"
)

// name[:TAG]
// return repository and tag
func Parse(input string) (string, string) {
	s := strings.Split(input, ":")
	if len(s) == 1 {
		return s[0], "latest"
	}
	if len(s) == 2 {
		return s[0], s[1]
	}
	return "", ""
}
