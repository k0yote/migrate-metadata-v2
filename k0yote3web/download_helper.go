package k0yote3web

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type downloadHelper struct {
	endpoints []string
}

func newDownloadHelper(opts *DownloadMetaOptions) (*downloadHelper, error) {
	var (
		endpoints = make([]string, 0)
		err       error
	)

	if opts != nil {
		endpoints, err = makeMetadataEndpointList(opts.BaseURL, opts.StartTokenID, opts.EndTokenID)
		if err != nil {
			return nil, err
		}
	}

	return &downloadHelper{
		endpoints: endpoints,
	}, nil
}

func (d *downloadHelper) updateEndpoints(endpoints []string) {
	if len(endpoints) > 0 {
		d.endpoints = endpoints
	}
}

func (d *downloadHelper) downloadMultipleFiles() ([]DownloadCh, error) {
	done := make(chan DownloadCh, len(d.endpoints))
	errch := make(chan error, len(d.endpoints))
	for _, endpoint := range d.endpoints {
		go func(endpoint string) {
			b, err := downloadFile(endpoint)
			if err != nil {
				errch <- err
				done <- DownloadCh{}
				return
			}

			done <- DownloadCh{
				Endpoint: endpoint,
				Data:     b,
			}
			errch <- nil
		}(endpoint)
	}

	downloadArr := make([]DownloadCh, 0)
	var errStr string
	for i := 0; i < len(d.endpoints); i++ {
		downloadArr = append(downloadArr, <-done)
		if err := <-errch; err != nil {
			errStr = errStr + " " + err.Error()
		}
	}

	var err error
	if errStr != "" {
		err = errors.New(errStr)
	}

	return downloadArr, err
}

func downloadFile(endpoint string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	req.Header.Add("Accept", "application/json")

	client := http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, errors.New(response.Status)
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func makeMetadataEndpointList(baseURL string, startTokenID int, endTokenID int) ([]string, error) {
	endpoints := []string{}

	for i := startTokenID; i <= endTokenID; i++ {
		endpoint, err := url.JoinPath(fmt.Sprintf(baseURL, i))
		if err != nil {
			return nil, err
		}

		endpoints = append(endpoints, endpoint)
	}

	return endpoints, nil
}
