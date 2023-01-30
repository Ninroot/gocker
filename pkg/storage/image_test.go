package storage

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/ninroot/gocker/pkg/testutil"
)

func TestCreateImage(t *testing.T) {
	root := testutil.MkdirTempTest(t, "")
	defer os.RemoveAll(root)

	storeDir := filepath.Join(root, "store")
	os.Mkdir(storeDir, 0755)
	store := NewImageStore(storeDir, testutil.Testfs{})

	imgDir := filepath.Join(root, "img")
	os.Mkdir(imgDir, 0755)
	_, err := os.OpenFile(filepath.Join(imgDir, "empty"), os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		log.Fatal("Could not make empty file", err)
	}

	r, w := io.Pipe()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		wg.Done()
		err := testutil.CreateTarball(w, []string{})
		if err != nil {
			log.Fatal("Could not create empty image tar file: ", err)
		}
		w.Close()
	}()

	h, err := store.CreateImage(r, testutil.RandID())
	wg.Wait()

	if err != nil {
		log.Fatal("Could not create image", err)
	}
	if h == nil {
		log.Fatal("Could not create image: CreateImage returned nil")
	}
}
