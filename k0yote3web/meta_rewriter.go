package k0yote3web

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type rewriteHelper struct {
	inputDir         string
	outputDir        string
	ipfsImageBaseURL string
}

func newRewriteHelper(ipfsImageBaseURL string) rewriteHelper {
	return rewriteHelper{
		inputDir:         metadataFolderName,
		outputDir:        uploadFolderName,
		ipfsImageBaseURL: ipfsImageBaseURL,
	}
}

func (r rewriteHelper) rewrite() error {
	metaDir, err := getSavePath(r.inputDir)
	if err != nil {
		return err
	}

	outputDir, err := getSavePath(r.outputDir)
	if err != nil {
		return err
	}

	metaFiles, err := os.ReadDir(metaDir)
	if err != nil {
		return err
	}

	for _, metaFile := range metaFiles {
		inputPath := filepath.Join(metaDir, metaFile.Name())

		b, err := os.ReadFile(inputPath)
		if err != nil {
			return err
		}

		m := MetaData{}
		if err := json.Unmarshal(b, &m); err != nil {
			return err
		}

		filename := filepath.Base(m.Image)
		newImagePath := filepath.Join(r.ipfsImageBaseURL, filename)
		m.Image = newImagePath

		metaByte, err := json.Marshal(m)
		if err != nil {
			return err
		}

		if err := saveJson(metaByte, outputDir, metaFile.Name()); err != nil {
			return err
		}
	}

	return nil
}
