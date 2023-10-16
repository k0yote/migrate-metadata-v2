package k0yote3web

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/ipfs/go-cid"
	ipfsapi "github.com/ipfs/kubo/client/rpc"

	coreiface "github.com/ipfs/boxo/coreiface"
	caopts "github.com/ipfs/boxo/coreiface/options"
	ipfsPath "github.com/ipfs/boxo/coreiface/path"
	ipfsFiles "github.com/ipfs/boxo/files"
)

type IpfsUploader struct {
	opts *IPFSOptions
}

func newIpfsUploader(opts *IPFSOptions) (*IpfsUploader, error) {
	if opts == nil {
		return nil, fmt.Errorf("provider type is required")
	}

	opts.ApiURL = defaultIpfsAPI
	if opts.ProviderType == IPFS_INFURA {
		opts.ApiURL = infuraAPI
	}

	return &IpfsUploader{
		opts: opts,
	}, nil
}

func (h *IpfsUploader) Upload(path string) (cid.Cid, error) {
	return h.upload(path)
}

func (h *IpfsUploader) GetGatewayUrl() string {
	if h.opts.ProviderType == IPFS_INFURA {
		return publicIpfsGatewayUrl
	}

	return defaultIpfsGatewayUrl
}

func (h *IpfsUploader) upload(path string) (cid.Cid, error) {
	var (
		cid        cid.Cid
		httpClient = &http.Client{}
	)

	client, err := ipfsapi.NewURLApiWithClient(h.opts.ApiURL, httpClient)
	if err != nil {
		return cid, err
	}

	if h.opts != nil && h.opts.ProviderType != IPFS_LOCAL {
		h.basicAuth(client)
	}

	stat, err := os.Lstat(path)
	if err != nil {
		return cid, err
	}

	file, err := ipfsFiles.NewSerialFile(path, false, stat)
	if err != nil {
		return cid, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	defer func() {
		signal.Stop(c)
		cancel()
	}()
	go func() {
		select {
		case <-c:
			cancel()
		case <-ctx.Done():
		}
	}()

	var res ipfsPath.Resolved
	errCh := make(chan error, 1)
	events := make(chan interface{}, 8)
	start := time.Now()

	go func() {
		var err error
		defer close(events)
		res, err = client.Unixfs().Add(ctx, file, caopts.Unixfs.Pin(h.opts.Pin), caopts.Unixfs.Progress(true), caopts.Unixfs.Events(events))
		errCh <- err
	}()

	for event := range events {
		output, ok := event.(*coreiface.AddEvent)
		if !ok {
			panic("unknown event type")
		}

		if output.Path != nil && output.Name != "" {
			if h.opts.Verbose {
				log.Printf("Added %v %v | Bytes: %v | Size: %v\n", output.Name, output.Path, output.Bytes, output.Size)
			} else {
				log.Printf("Added %v\n", output.Name)
			}
		}
	}

	if err := <-errCh; err != nil {
		elapse(start)
		return cid, err
	}

	elapse(start)
	return res.Cid(), nil
}

func (h *IpfsUploader) basicAuth(client *ipfsapi.HttpApi) {
	basicAuth := fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(h.opts.ProjectID+":"+h.opts.Secret)))
	client.Headers.Add("Authorization", basicAuth)
}

func elapse(start time.Time) {
	duration := time.Since(start)
	log.Println("elapsed time: ", duration)
}
