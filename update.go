package goup

import (
	"net/url"
	"time"
)

type Update struct {
	URL         url.URL
	Version     string
	ChecksumURL url.URL
	Date        time.Time
	Size        uint64
	OS          string
	Arch        string
	Extras      any
}
