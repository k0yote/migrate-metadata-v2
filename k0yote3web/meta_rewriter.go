package k0yote3web

import (
	"encoding/json"
	"os"
	"path/filepath"
)

var counter int

type MetaRewriter struct {
	inputDir         string
	outputDir        string
	ipfsImageBaseURL string
}

func newMetaRewriter(ipfsImageBaseURL, inputDir, outputDir string) (*MetaRewriter, error) {

	in := metadataFolderName
	out := uploadFolderName

	if len(inputDir) > 0 {
		in = inputDir
		if _, err := getSavePath(in); err != nil {
			return nil, err
		}
	}

	if len(outputDir) > 0 {
		out = outputDir
		if _, err := getSavePath(out); err != nil {
			return nil, err
		}
	}

	counter = 0

	return &MetaRewriter{
		inputDir:         in,
		outputDir:        out,
		ipfsImageBaseURL: ipfsImageBaseURL,
	}, nil
}

func (r *MetaRewriter) Rewrite() error {
	return r.rewrite()
}

func (r MetaRewriter) rewrite() error {
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

		counter++
	}

	return nil
}

func (r *MetaRewriter) Counter() int {
	return counter
}
