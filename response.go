package dose

import (
	"time"
)

type AddedResponse struct {
	Path string
}

func (r AddedResponse) MessageType() MessageType {
	return AddedResponseMessage
}

type CanceledResponse struct {
	Path string
}

func (r CanceledResponse) MessageType() MessageType {
	return CanceledResponseMessage
}

type InfoResponse struct {
	Download DownloadResponse
}

func (r InfoResponse) MessageType() MessageType {
	return InfoResponseMessage
}

type ServerInfoResponse struct {
	Downloads []DownloadResponse
}

func (r ServerInfoResponse) MessageType() MessageType {
	return ServerInfoResponseMessage
}

type ErrorResponse struct {
	Error string
}

func (r ErrorResponse) MessageType() MessageType {
	return ErrorResponseMessage
}

type DownloadResponse struct {
	Url       string
	Path      string
	Status    DownloadStatus
	BytesRead uint
	Filesize  *uint
	StartTime time.Time
}

func (r DownloadResponse) MessageType() MessageType {
	return DownloadResponseMessage
}

type AuthResponse struct{}

func (r AuthResponse) MessageType() MessageType {
	return AuthResponseMessage
}
