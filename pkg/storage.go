package pkg

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/containerd/btrfs"
)

type ImageStore struct {
	rootDir string
}

func NewImageStore(rootDir string) ImageStore {
	return ImageStore{
		rootDir: rootDir,
	}
}

func (s ImageStore) CreateImage(reader io.ReadCloser, image ImageId) error {
	if err := os.MkdirAll(s.rootDir, 0700); err != nil {
		return err
	}

	rootDir := filepath.Join(s.rootDir, image.Digest)
	prs, err := Exist(rootDir)
	if err != nil {
		return err
	}
	if prs {
		log.Printf("image <%s> already exists", rootDir)
		return nil
	}

	if err := btrfs.SubvolCreate(rootDir); err != nil {
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

func (s ImageStore) RemoveImage(img *ImageId) (*ImageId, error) {
	found, err := s.FindImage(img)
	if err != nil {
		return nil, err
	}
	if found == nil {
		return nil, fmt.Errorf("image <%s:%s> not found", img.Name, img.Tag)
	}
	return found, os.RemoveAll(path.Join(s.rootDir, found.Digest))
}

func (s ImageStore) ListImages() ([]*ImageId, error) {
	return s.findImages(nil)
}

// FindImage returns the image found in the store or nil if not found.
func (s ImageStore) FindImage(image *ImageId) (*ImageId, error) {
	found, err := s.findImages(image)
	if err != nil {
		return nil, err
	}
	if len(found) == 0 {
		return nil, nil
	}
	return found[0], nil
}

func (s ImageStore) findImages(image *ImageId) ([]*ImageId, error) {
	items, err := os.ReadDir(s.rootDir)
	if err != nil {
		return nil, err
	}
	images := make([]*ImageId, 0)
	for _, item := range items {
		if item.IsDir() {
			source := path.Join(s.rootDir, item.Name(), "source")
			ok, err := Exist(source)
			if err != nil {
				return nil, err
			}
			if ok {
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
				if image != nil && img.Name == image.Name && img.Tag == image.Tag {
					return []*ImageId{&img}, nil
				}
				images = append(images, &img)
			}
		}
	}
	return images, nil
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

	toSymLink := map[string]string{}

L:
	for {
		header, err := tr.Next()

		switch {

		// if no more files are found return
		case err == io.EOF:
			break L

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
		case tar.TypeSymlink:
			// save symlink to be done later when files will exist
			toSymLink[target] = header.Linkname
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

	for target, linkName := range toSymLink {
		// link target -> linkName
		err := os.Symlink(linkName, target)
		if err != nil {
			return err
		}
	}

	return nil
}
