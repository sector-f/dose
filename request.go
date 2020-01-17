package dose

import (
	"fmt"
)

type AddRequest struct {
	Url  string `json:"url"`
	Path string `json:"path"`
}

func (r AddRequest) MessageType() MessageType {
	return AddRequestMessage
}

func (r AddRequest) String() string {
	return fmt.Sprintf("AddRequest: %s %s", r.Url, r.Path)
}

type CancelRequest struct {
	Path string `json:"path"`
}

func (r CancelRequest) MessageType() MessageType {
	return CancelRequestMessage
}

func (r CancelRequest) String() string {
	return fmt.Sprintf("CancelRequest: %s", r.Path)
}

type RemoveRequest struct {
	Path string `json:"path"`
}

func (r RemoveRequest) MessageType() MessageType {
	return RemoveRequestMessage
}

func (r RemoveRequest) String() string {
	return fmt.Sprintf("RemoveRequest: %s", r.Path)
}

type InfoRequest struct {
	Path string `json:"path"`
}

func (r InfoRequest) MessageType() MessageType {
	return InfoRequestMessage
}

func (r InfoRequest) String() string {
	return fmt.Sprintf("InfoRequest: %s", r.Path)
}

type ServerInfoRequest struct{}

func (r ServerInfoRequest) MessageType() MessageType {
	return ServerInfoRequestMessage
}

func (r ServerInfoRequest) String() string {
	return fmt.Sprintf("ServerInfoRequest")
}

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (r AuthRequest) MessageType() MessageType {
	return AuthRequestMessage
}

func (r AuthRequest) String() string {
	return fmt.Sprintf("AuthRequest")
}
