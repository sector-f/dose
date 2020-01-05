package dose

type DownloadStatus uint

const (
	Queued     DownloadStatus = 0
	InProgress DownloadStatus = 1
	Completed  DownloadStatus = 2
	Failed     DownloadStatus = 3
	Canceled   DownloadStatus = 4
	Paused     DownloadStatus = 5
)

func (s DownloadStatus) String() string {
	switch s {
	case Queued:
		return "Queued"
	case InProgress:
		return "InProgress"
	case Completed:
		return "Completed"
	case Failed:
		return "Failed"
	case Canceled:
		return "Canceled"
	case Paused:
		return "Paused"
	default:
		return "Undefined"
	}
}
