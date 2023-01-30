package storage

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"fmt"
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
		err := CreateTarball(w, []string{})
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

func testPipe(w io.Writer) {
	fmt.Fprint(w, "Hello")
}

// Credit: https://gist.github.com/maximilien/328c9ac19ab0a158a8df
func CreateTarball(w io.Writer, filePaths []string) error {
	gzipWriter := gzip.NewWriter(w)
	// defer gzipWriter.Close()
	defer func() {
		if err := gzipWriter.Close(); err != nil {
			log.Fatal("Failed to close", err)
		}
	}()

	tarWriter := tar.NewWriter(gzipWriter)
	// defer tarWriter.Close()
	defer func() {
		if err := tarWriter.Close(); err != nil {
			log.Fatal("Failed to close", err)
		}
	}()

	for _, filePath := range filePaths {
		err := addFileToTarWriter(filePath, tarWriter)
		if err != nil {
			return errors.New(fmt.Sprintf("Could not add file '%s', to tarball, got error '%s'", filePath, err.Error()))
		}
	}

	return nil
}

func addFileToTarWriter(filePath string, writer *tar.Writer) error {
	file, err := os.Open(filePath)
	if err != nil {
		return errors.New(fmt.Sprintf("Could not open file '%s', got error '%s'", filePath, err.Error()))
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return errors.New(fmt.Sprintf("Could not get stat for file '%s', got error '%s'", filePath, err.Error()))
	}

	header := &tar.Header{
		Name:    filePath,
		Size:    stat.Size(),
		Mode:    int64(stat.Mode()),
		ModTime: stat.ModTime(),
	}

	err = writer.WriteHeader(header)
	if err != nil {
		return errors.New(fmt.Sprintf("Could not write header for file '%s', got error '%s'", filePath, err.Error()))
	}

	_, err = io.Copy(writer, file)
	if err != nil {
		return errors.New(fmt.Sprintf("Could not copy the file '%s' data to the tarball, got error '%s'", filePath, err.Error()))
	}

	return nil
}
