package k0yote3web

import (
	"log"
	"math"
	"time"

	"golang.org/x/exp/slices"
)

type Download struct {
	downloadHelper *downloadHelper
	imgHelper      *imageHelper
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
		downloadHelper: helper,
		imgHelper:      imgHelper,
	}, nil
}

func (d Download) DownloadAndSaveMetadata() error {

	savePath, err := getSavePath(metadataFolderName)
	if err != nil {
		return err
	}

	maxPage := d.GetMetaMaxPage()

	downloadAndSavedCount := 0
	for i := 1; i <= maxPage; i++ {
		endpoints, _, _ := pagination(d.downloadHelper.endpoints, i, fetchDownloadMetaLimit)

		downloadList, err := downloadMultipleFiles(endpoints)
		if err != nil {
			return err
		}

		for _, download := range downloadList {
			if err := saveJson(download.Data, savePath, download.Endpoint); err != nil {
				return err
			}
		}

		downloadAndSavedCount += len(downloadList)
		log.Println("downloaded and saved count: ", downloadAndSavedCount)
		time.Sleep(waitTime)
	}

	return nil
}

func (d Download) DownloadAndSaveImage() error {
	endpoints, err := d.imgHelper.getImageURLByMetadata()
	if err != nil {
		return err
	}

	slices.Sort(endpoints)
	uniqueEndpoints := slices.Compact(endpoints)

	d.imgHelper.updateEndpoints(uniqueEndpoints)

	maxPage := d.GetImageMaxPage()
	savePath, err := getSavePath(imageFolderName)
	if err != nil {
		return err
	}

	downloadAndSavedCount := 0
	for i := 1; i <= maxPage; i++ {
		endpoints, _, _ := pagination(d.imgHelper.endpoints, i, fetchDownloadImageLimit)

		downloadList, err := downloadMultipleFiles(endpoints)
		if err != nil {
			return err
		}

		for _, download := range downloadList {
			if err := saveImage(download.Data, savePath, download.Endpoint); err != nil {
				return err
			}
		}
		downloadAndSavedCount += len(downloadList)
		log.Println("downloaded and saved count: ", downloadAndSavedCount)
		time.Sleep(waitTime)
	}

	return nil
}

func (d Download) GetMetaMaxPage() int {
	_, _, maxPage := pagination(d.downloadHelper.endpoints, 1, fetchDownloadMetaLimit)
	return maxPage
}

func (d Download) GetImageMaxPage() int {
	_, _, maxPage := pagination(d.imgHelper.endpoints, 1, fetchDownloadImageLimit)
	return maxPage
}

func (d Download) GetDownloadMetaCount() int {
	return len(d.downloadHelper.endpoints)
}

func (d Download) GetDownloadImageCount() int {
	return len(d.imgHelper.endpoints)
}

func pagination[T comparable](x []T, page int, perPage int) (data []T, currentPage int, lastPage int) {
	lastPage = int(math.Ceil(float64(len(x)) / float64(perPage)))
	currentPage = page

	if page < 1 {
		page = 1
	} else if lastPage < page {
		page = lastPage
	}

	if page == lastPage {
		data = x[(page-1)*perPage:]
	} else {
		data = x[(page-1)*perPage : page*perPage]
	}

	return
}
