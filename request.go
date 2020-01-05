package dose

type AddRequest struct {
	Url  string `json:"url"`
	Path string `json:"path"`
}

type CancelRequest struct {
	Path string `json:"path"`
}

type RemoveRequest struct {
	Path string `json:"path"`
}

type InfoRequest struct {
	Path string `json:"path"`
}

type ServerInfoRequest struct{}
