package pkg

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
)

type ImageStore struct {
	rootDir string
}

func NewImageStore(rootDir string) ImageStore {
	return ImageStore{
		rootDir: rootDir,
	}
}

func (store ImageStore) CreateImage(reader io.ReadCloser, image ImageId) error {
	rootDir := filepath.Join(store.rootDir, image.Digest)
	prs, err := Exist(rootDir)
	if err != nil {
		return err
	}
	if prs {
		log.Printf("image <%s> already exists", rootDir)
		return nil
	}

	if err := os.MkdirAll(rootDir, 0700); err != nil {
		return err
	}

	if err := Untar(rootDir+"/rootfs", reader); err != nil {
		return err
	}

	source, err := os.Create(filepath.Join(rootDir, "source"))
	if err != nil {
		return err
	}
	encoder := json.NewEncoder(source)
	if err := encoder.Encode(image); err != nil {
		return err
	}

	log.Printf("image <%s> stored at <%s>", image.Digest, rootDir)
	return nil
}

func (store ImageStore) FindImage(image ImageId) (*ImageId, error) {
	items, err := os.ReadDir(store.rootDir)
	if err != nil {
		return nil, err
	}
	for _, item := range items {
		if item.IsDir() {
			source := path.Join(store.rootDir, item.Name(), "source")
			prs, err := Exist(source)
			log.Println("source", source)
			if err != nil {
				return nil, err
			}
			if prs {
				file, err := os.Open(source)
				if err != nil {
					return nil, err
				}
				defer file.Close()
				content, err := io.ReadAll(file)
				if err != nil {
					return nil, err
				}
				var img ImageId
				if err := json.Unmarshal(content, &img); err != nil {
					return nil, err
				}
				if img.Name == image.Name && img.Tag == image.Tag {
					return &img, nil
				}
			}
		}
	}
	return nil, nil
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
