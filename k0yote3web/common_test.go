package k0yote3web

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSavePath(t *testing.T) {
	got, err := getSavePath(metadataFolderName)
	assert.NoError(t, err)

	want := filepath.Join(rootDir(), saveFolderName, metadataFolderName)
	assert.Equal(t, want, got)
}

func TestRootDir(t *testing.T) {
	currentPath, err := os.Getwd()
	assert.NoError(t, err)

	rootPath := rootDir()
	assert.NotEqual(t, currentPath, rootPath)
}

func TestGetFilename(t *testing.T) {
	endpoint := "http://example.com/abc/def/1"

	got, err := getFilename(endpoint)
	assert.NoError(t, err)

	assert.Equal(t, "1", got)
}

func TestSaveJson(t *testing.T) {

}

func TestSaveImage(t *testing.T) {

}
