package testutil

import (
	"log"
	"os"
	"testing"
)

func MkdirTempTest(t *testing.T, pattern string) string {
	d, err := os.MkdirTemp("", t.Name()+pattern)
	if err != nil {
		log.Fatal(err)
	}
	return d
}
