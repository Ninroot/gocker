package container

import (
	"log"
	"strings"

	"github.com/google/uuid"
)

func RandID() string {
	id, err := uuid.NewRandom()
	if err != nil {
		log.Fatal("Enable to generate new UUID: ", err)
	}
	return strings.ReplaceAll(id.String(), "-", "")
}
