package storage

import (
	"log"
	"os"
	"testing"

	"github.com/ninroot/gocker/pkg/container"
	"github.com/ninroot/gocker/pkg/testutil"
)

func TestCreateContainer(t *testing.T) {
	img := testutil.MkdirTempTest(t, "img")
	defer os.RemoveAll(img)
	con := testutil.MkdirTempTest(t, "con")
	defer os.RemoveAll(con)

	s := NewContainerStore(img, testutil.Testfs{})
	_, err := s.CreateContainer(container.RandID(), con)
	if err != nil {
		log.Fatal("ContainerStore failed: CreateContainer")
	}
}
