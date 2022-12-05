package pkg

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
