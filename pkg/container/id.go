package container

import (
	"strings"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func RandID() string {
	id, err := uuid.NewRandom()
	if err != nil {
		logrus.Fatal("Enable to generate new UUID: ", err)
	}
	return strings.ReplaceAll(id.String(), "-", "")
}
