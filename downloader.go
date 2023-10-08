package goup

import (
	"github.com/go-resty/resty/v2"
	"io"
)

type Downloader interface {
	Download() (io.ReadCloser, error)
}

type restyDownloader struct {
	client *resty.Client
	url    string
}

func Download(url string) Downloader {
	return &restyDownloader{client: resty.New(), url: url}
}

func DownloadWithResty(url string, client *resty.Client) Downloader {
	return &restyDownloader{client: client, url: url}
}

func (r *restyDownloader) Download() (io.ReadCloser, error) {
	resp, err := r.client.R().
		SetDoNotParseResponse(true).
		SetHeader("Accept", "application/octet-stream").
		Get(r.url)
	if err != nil {
		return nil, err
	}
	return resp.RawBody(), nil
}
