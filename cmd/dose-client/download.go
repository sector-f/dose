package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/sector-f/dose"
)

func download(conn net.Conn, url, filepath string) {
	data := dose.AddRequest{
		url,
		filepath,
	}

	conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
	_, err := dose.WriteMessage(conn, data)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error sending message to server")
		return
	}

	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	response, err := dose.ReadMessage(conn)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error receiving message from server")
		return
	}

	switch r := response.(type) {
	case *dose.AddedResponse:
		fmt.Printf("Queued download to %v\n", r.Path)
	case *dose.ErrorResponse:
		fmt.Printf("Received error: %s\n", r.Error)
	default:
		fmt.Printf("Received unexpected message of type %v\n", r.MessageType())
	}
}
