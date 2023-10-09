package k0yote3web

import (
	"bytes"
	"encoding/json"
	"log"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func getSavePath(folderName string) (string, error) {
	wd := rootDir()
	tmpDir := path.Join(wd, folderName)
	if f, err := os.Stat(tmpDir); os.IsNotExist(err) || !f.IsDir() {
		if err := os.Mkdir(tmpDir, 0777); err != nil {
			return "", err
		}
	}

	return tmpDir, nil
}

func rootDir() string {
	rootDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	for {
		_, err := os.ReadFile(filepath.Join(rootDir, "go.mod"))
		if os.IsNotExist(err) {
			if rootDir == filepath.Dir(rootDir) {
				// at the root
				break
			}
			rootDir = filepath.Dir(rootDir)
			continue
		} else if err != nil {
			log.Fatal(err)
		}
		break
	}

	return rootDir
}

func getFilename(endpoint string) (string, error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		return "", err
	}

	path := u.Path
	segments := strings.Split(path, "/")

	return segments[len(segments)-1], nil
}

func saveJson(data []byte, savePath, endpoint string) error {

	filename, err := getFilename(endpoint)
	if err != nil {
		return err
	}

	jsonObj := &MetaData{}
	if err := json.Unmarshal(data, &jsonObj); err != nil {
		return err
	}

	outputPath := path.Join(savePath, filename)

	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}

	defer file.Close()
	if err := json.NewEncoder(file).Encode(jsonObj); err != nil {
		return err
	}
	return nil
}

func saveImage(data []byte, savePath, endpoint string) error {
	filename, err := getFilename(endpoint)
	if err != nil {
		return err
	}

	outputPath := path.Join(savePath, filename)

	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.ReadFrom(bytes.NewReader(data))
	if err != nil {
		return err
	}

	return nil
}
