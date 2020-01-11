package dose

type AddRequest struct {
	Url  string `json:"url"`
	Path string `json:"path"`
}

func (r AddRequest) MessageType() MessageType {
	return AddRequestMessage
}

type CancelRequest struct {
	Path string `json:"path"`
}

func (r CancelRequest) MessageType() MessageType {
	return CancelRequestMessage
}

type RemoveRequest struct {
	Path string `json:"path"`
}

func (r RemoveRequest) MessageType() MessageType {
	return RemoveRequestMessage
}

type InfoRequest struct {
	Path string `json:"path"`
}

func (r InfoRequest) MessageType() MessageType {
	return InfoRequestMessage
}

type ServerInfoRequest struct{}

func (r ServerInfoRequest) MessageType() MessageType {
	return ServerInfoRequestMessage
}

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (r AuthRequest) MessageType() MessageType {
	return AuthRequestMessage
}
