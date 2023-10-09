package k0yote3web

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type imageHelper struct {
	dir string
}

func newImageHelper() (*imageHelper, error) {
	dir, err := getSavePath(imageFolderName)
	if err != nil {
		return nil, err
	}

	return &imageHelper{
		dir: dir,
	}, nil
}

func (s *imageHelper) getImageURLByMetadata() ([]string, error) {
	files, err := os.ReadDir(s.dir)
	if err != nil {
		return nil, err
	}

	endpoints := []string{}
	for _, file := range files {
		filename := filepath.Join(s.dir, file.Name())
		b, err := os.ReadFile(filename)
		if err != nil {
			return nil, err
		}

		var meta MetaData
		err = json.Unmarshal(b, &meta)
		if err != nil {
			return nil, err
		}

		endpoints = append(endpoints, meta.Image)
	}

	return endpoints, nil
}
