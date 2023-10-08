package goup

import (
	"time"
)

type Update struct {
	GetFile  Downloader
	Version  string
	Checksum []byte
	Time     time.Time
	Size     uint64
	OS       string
	Arch     string
	Extras   any
}
