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
					default:
						return
					}
				}(conn)
			}
		}(*l)
	}

	wg.Wait()
}

type Download struct {
	Url       string
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
