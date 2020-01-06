package dose

import (
	"time"
)

type AddedResponse struct {
	Path string
}

type CanceledResponse struct {
	Path string
}

type InfoResponse struct {
	Download DownloadResponse
}

type ServerInfoResponse struct {
	Downloads []DownloadResponse
}

type ErrorResponse struct {
	error string
}

type DownloadResponse struct {
	Url       string
	Path      string
	Status    DownloadStatus
	BytesRead uint
	Filesize  *uint
	StartTime time.Time
}
