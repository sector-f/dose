package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/sector-f/dose"
)

type DownloadServer struct {
	Downloads map[string]*Download // Map path to Download
}

type Download struct {
	Url       string
	Status    dose.DownloadStatus
	BytesRead uint
	Filesize  *uint
	StartTime time.Time
	Cancel    context.CancelFunc
}

type readerFunc func(p []byte) (n int, err error)

func (rf readerFunc) Read(p []byte) (n int, err error) { return rf(p) }

func (s *DownloadServer) Download(url string, path string) {
	// TODO: Why did I put this here???
	// _, prs := s.Downloads[filepath.Clean(path)]
	// if !prs {
	// 	return errors.New("Download not found")
	// }

	download := &Download{
		Url:        url,
		Status:     dose.Queued,
		BytesRead:  0,
		BytesTotal: nil,
		StartTime:  time.Now(),
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
		dl.Cancel()
		dl.Status = dose.Canceled
		return nil
	}

	return errors.New("Download not found")
}

func main() {
	downloadServer := DownloadServer{make(map[string]*Download)}

	listener, err := net.Listen("unix", "/tmp/dose.socket")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for {
		conn, _ := listener.Accept()
		go func(c net.Conn) {
			defer c.Close()

			var headerBytes [4]byte
			_, err := c.Read(headerBytes[:])
			if err != nil {
				log.Println(err)
				return
			}

			header := dose.ParseHeader(headerBytes)

			buf := make([]byte, header.Length)
			c.Read(buf)

			request, err := dose.ParseBody(header.MessageType, buf)
			if err != nil {
				fmt.Println(err)
				return
			}

			switch r := request.(type) {
			case dose.AddRequest:
				log.Printf("AddRequest: %s\t%s\n", r.Url, r.Path)
				downloadServer.Download(r.Url, r.Path)
			case dose.CancelRequest:
				log.Printf("CancelRequest: %s\n", r.Path)
				downloadServer.Cancel(r.Path)
			default:
				return
			}
		}(conn)
	}
}
