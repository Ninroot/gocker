package pkg

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/ninroot/gocker/config"
)

func CreateImage(reader io.ReadCloser, digest string) error {
	rootDir := config.DefaultImageStoreRootDir + digest
	prs, err := Exist(rootDir)
	if err != nil {
		return err
	}
	if prs {
		return fmt.Errorf("image <%s> already exists", rootDir)
	}

	if err := os.MkdirAll(rootDir, 0700); err != nil {
		return err
	}

	if err := Untar(rootDir+"/rootfs", reader); err != nil {
		return err
	}

	log.Printf("image <%s> stored at <%s>", digest, rootDir)
	return nil
}

// Untar takes a destination path and a reader; a tar reader loops over the tarfile
// creating the file structure at 'dst' along the way, and writing any files
// credit: https://gist.github.com/sdomino/635a5ed4f32c93aad131/
func Untar(dst string, r io.Reader) error {

	gzr, err := gzip.NewReader(r)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()

		switch {

		// if no more files are found return
		case err == io.EOF:
			return nil

		// return any other error
		case err != nil:
			return err

		// if the header is nil, just skip it (not sure how this happens)
		case header == nil:
			continue
		}

		// the target location where the dir/file should be created
		target := filepath.Join(dst, header.Name)

		// the following switch could also be done using fi.Mode(), not sure if there
		// a benefit of using one vs. the other.
		// fi := header.FileInfo()

		// check the file type
		switch header.Typeflag {

		// if its a dir and it doesn't exist create it
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0755); err != nil {
					return err
				}
			}

		// if it's a file create it
		case tar.TypeReg:
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}

			// copy over contents
			if _, err := io.Copy(f, tr); err != nil {
				return err
			}

			// manually close here after each file operation; defering would cause each file close
			// to wait until all operations have completed.
			f.Close()
		}
	}
}
