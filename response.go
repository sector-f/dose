package dose

import (
	"fmt"
	"time"
)

type AddedResponse struct {
	Url  string `json:"url"`
	Path string
}

func (r AddedResponse) MessageType() MessageType {
	return AddedResponseMessage
}

func (r AddedResponse) String() string {
	return fmt.Sprintf("AddedResponse: %s %s", r.Url, r.Path)
}

type CanceledResponse struct {
	Path string
}

func (r CanceledResponse) MessageType() MessageType {
	return CanceledResponseMessage
}

func (r CanceledResponse) String() string {
	return fmt.Sprintf("CanceledResponse: %s", r.Path)
}

type InfoResponse struct {
	Download DownloadResponse
}

func (r InfoResponse) MessageType() MessageType {
	return InfoResponseMessage
}

func (r InfoResponse) String() string {
	return fmt.Sprintf("InfoResponse")
}

type ServerInfoResponse struct {
	Downloads []DownloadResponse
}

func (r ServerInfoResponse) MessageType() MessageType {
	return ServerInfoResponseMessage
}

func (r ServerInfoResponse) String() string {
	return fmt.Sprintf("ServerInfoResponse")
}

type ErrorResponse struct {
	Error string
}

func (r ErrorResponse) MessageType() MessageType {
	return ErrorResponseMessage
}

func (r ErrorResponse) String() string {
	return fmt.Sprintf("ErrorResponse: %v", r.Error)
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

func (r DownloadResponse) String() string {
	return fmt.Sprintf("DownloadResponse: %s", r.Path)
}

type AuthResponse struct{}

func (r AuthResponse) MessageType() MessageType {
	return AuthResponseMessage
}

func (r AuthResponse) String() string {
	return fmt.Sprintf("AuthResponse")
}
