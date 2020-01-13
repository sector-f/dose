package main

import (
	"context"
	"errors"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/sector-f/dose"
)

type DownloadServer struct {
	Downloads map[string]*Download // Map path to Download
}

func runDownloadServer(listeners []*net.Listener) {
	downloadServer := DownloadServer{make(map[string]*Download)}

	var wg sync.WaitGroup
	for _, l := range listeners {
		wg.Add(1)
		log.Printf("Listening on %s\n", (*l).Addr().String())
		go func(l net.Listener) {
			defer wg.Done()
			for {
				conn, _ := l.Accept()
				go func(c net.Conn) {
					defer c.Close()

					request, err := dose.ReadMessage(c)
					if err != nil {
						log.Println(err)
						return
					}

					switch r := request.(type) {
					case *dose.AddRequest:
						log.Printf("AddRequest: %s\t%s\n", r.Url, r.Path)
						downloadServer.Download(r.Url, r.Path)
						dose.WriteMessage(c, dose.AddedResponse{r.Path})
					case *dose.CancelRequest:
						log.Printf("CancelRequest: %s\n", r.Path)
						err := downloadServer.Cancel(r.Path)
						if err != nil {
							dose.WriteMessage(c, dose.ErrorResponse{err.Error()})
						} else {
							dose.WriteMessage(c, dose.CanceledResponse{r.Path})
						}
					case *dose.ServerInfoRequest:
						log.Printf("ServerInfoRequest")
						dose.WriteMessage(c, dose.ServerInfoResponse{downloadServer.ServerInfo()})
					default:
						dose.WriteMessage(c, dose.ErrorResponse{"Unimplemented function"})
					}
				}(conn)
			}
		}(*l)
	}

	wg.Wait()
}

type Download struct {
	Url       string
	Path      string
	Status    dose.DownloadStatus
	BytesRead uint
	Filesize  *uint
	StartTime time.Time
	Cancel    context.CancelFunc
	Mu        sync.Mutex
}

type readerFunc func(p []byte) (n int, err error)

func (rf readerFunc) Read(p []byte) (n int, err error) { return rf(p) }

func (s *DownloadServer) Download(url string, path string) {
	var mu sync.Mutex
	download := &Download{
		Url:       url,
		Path:      path,
		Status:    dose.Queued,
		BytesRead: 0,
		Filesize:  nil,
		StartTime: time.Now(),
		Mu:        mu,
	}
	s.Downloads[filepath.Clean(path)] = download

	go func() {
		out, err := os.Create(path)
		if err != nil {
			(*download).Status = dose.Failed
			return
		}
		defer out.Close()

		resp, err := http.Get(url)
		if err != nil {
			(*download).Status = dose.Failed
			return
		}
		defer resp.Body.Close()

		if len := resp.Header.Get("Content-Length"); len != "" {
			if asInt, err := strconv.Atoi(len); err == nil {
				asUint := uint(asInt)
				(*download).Filesize = &asUint
			}
		}

		newContext, cancelFunc := context.WithCancel(resp.Request.Context())
		download.Cancel = cancelFunc

		func(ctx context.Context, dst io.Writer, src io.Reader) {
			io.Copy(dst, readerFunc(func(p []byte) (int, error) {
				select {
				case <-ctx.Done():
					return 0, ctx.Err()
				default:
					read, err := src.Read(p)
					if err == nil {
						download.BytesRead += uint(read)
					}
					return read, err
				}
			}))
		}(newContext, out, resp.Body)
	}()
}

func (s *DownloadServer) Cancel(path string) error {
	dl, prs := s.Downloads[filepath.Clean(path)]
	if prs {
		dl.Mu.Lock()
		defer dl.Mu.Unlock()

		switch dl.Status {
		case dose.Queued, dose.InProgress, dose.Paused:
			dl.Cancel()
			dl.Status = dose.Canceled
			return nil
		case dose.Canceled:
			return errors.New("Download has already been canceled")
		case dose.Completed:
			return errors.New("Download has already completed")
		case dose.Failed:
			return errors.New("Attempted to cancel failed download")
		}
	}

	return errors.New("Download not found")
}

func (s *DownloadServer) ServerInfo() []dose.DownloadResponse {
	sorted := downloads([]dose.DownloadResponse{})
	for _, download := range s.Downloads {
		response := dose.DownloadResponse{
			Url:       download.Url,
			Path:      download.Path,
			Status:    download.Status,
			BytesRead: download.BytesRead,
			Filesize:  download.Filesize,
			StartTime: download.StartTime,
		}
		sorted = append(sorted, response)
	}
	sort.Sort(sorted)

	return sorted
}

type downloads []dose.DownloadResponse

func (d downloads) Len() int {
	return len(d)
}

func (d downloads) Less(i, j int) bool {
	return d[i].StartTime.Before(d[j].StartTime)
}

func (d downloads) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

//func main() {
//	var reviews_data_map = make(map[string]reviews_data)
//	reviews_data_map["1"] = reviews_data{date: time.Now().Add(12 * time.Hour)}
//	reviews_data_map["2"] = reviews_data{date: time.Now()}
//	reviews_data_map["3"] = reviews_data{date: time.Now().Add(24 * time.Hour)}
//	//Sort the map by date
//	date_sorted_reviews := make(timeSlice, 0, len(reviews_data_map))
//	for _, d := range reviews_data_map {
//		date_sorted_reviews = append(date_sorted_reviews, d)
//	}
//	fmt.Println(date_sorted_reviews)
//	sort.Sort(date_sorted_reviews)
//	fmt.Println(date_sorted_reviews)
//}
