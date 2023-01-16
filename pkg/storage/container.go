package storage

// func (s ImageHandle) CreateContainer(img *ImageId) (string, error) {
// 	found, err := s.findImages(img)
// 	if err != nil {
// 		return "", err
// 	}
// 	if found == nil {
// 		return "", fmt.Errorf("image <%s:%s> not found", img.Name, img.Digest)
// 	}

// 	contDir := filepath.Join(s.imageDir, "con")
// 	if err := os.MkdirAll(contDir, 0700); err != nil {
// 		return "", err
// 	}

// 	id := container.RandID()
// 	err = btrfs.SubvolSnapshot(filepath.Join(contDir, id), filepath.Join(s.imageDir, "img", img.Digest), false)
// 	if err != nil {
// 		return "", err
// 	}
// 	return id, nil
// }
