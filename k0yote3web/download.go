package k0yote3web

import (
	"path/filepath"
)

type Download struct {
	helper    *downloadHelper
	imgHelper *imageHelper
}

func newDownload(opts *DownloadMetaOptions) (*Download, error) {

	helper, err := newDownloadHelper(opts)
	if err != nil {
		return nil, err
	}

	imgHelper, err := newImageHelper()
	if err != nil {
		return nil, err
	}

	return &Download{
		helper:    helper,
		imgHelper: imgHelper,
	}, nil
}

func (d Download) DownloadAndSaveMetadata() error {
	downloadList, err := d.helper.downloadMultipleFiles()
	if err != nil {
		return err
	}

	savePath, err := getSavePath(filepath.Join(saveFolderName, metadataFolderName))
	if err != nil {
		return err
	}

	for _, download := range downloadList {
		if err := saveJson(download.Data, savePath, download.Endpoint); err != nil {
			return err
		}
	}

	return nil
}

func (d Download) DownloadAndSaveImage() error {
	endpoints, err := d.imgHelper.getImageURLByMetadata()
	if err != nil {
		return err
	}

	d.helper.updateEndpoints(endpoints)

	downloadList, err := d.helper.downloadMultipleFiles()
	if err != nil {
		return err
	}

	savePath, err := getSavePath(filepath.Join(saveFolderName, metadataFolderName))
	if err != nil {
		return err
	}

	for _, download := range downloadList {
		if err := saveImage(download.Data, savePath, download.Endpoint); err != nil {
			return err
		}
	}

	return nil
}
