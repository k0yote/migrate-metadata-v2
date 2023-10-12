package k0yote3web

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type imageHelper struct {
	metaDir   string
	endpoints []string
}

func newImageHelper() (*imageHelper, error) {
	metaDir, err := getSavePath(metadataFolderName)
	if err != nil {
		return nil, err
	}

	return &imageHelper{
		metaDir:   metaDir,
		endpoints: make([]string, 0),
	}, nil
}

func (s *imageHelper) updateEndpoints(endpoints []string) {
	if len(endpoints) > 0 {
		s.endpoints = endpoints
	}
}

func (s *imageHelper) getImageURLByMetadata() ([]string, error) {
	files, err := os.ReadDir(s.metaDir)
	if err != nil {
		return nil, err
	}

	endpoints := []string{}
	for _, file := range files {
		filename := filepath.Join(s.metaDir, file.Name())
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
