package container

import (
	"time"

	"github.com/ninroot/gocker/pkg/image"
)

type Container struct {
	ID      string      `json:"id"`
	Name    string      `json:"name"`
	Image   image.Image `json:"image"`
	Command string      `json:"command"`
	Args    []string    `json:"args"`

	CreatedAt time.Time `json:"createdAt"`
}
