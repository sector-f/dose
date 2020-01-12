package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/sector-f/dose"
)

func cancel(conn net.Conn, filepath string) {
	data := dose.CancelRequest{filepath}

	conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
	_, err := dose.WriteMessage(conn, data)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	response, err := dose.ReadMessage(conn)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	switch r := response.(type) {
	case *dose.CanceledResponse:
		fmt.Printf("Canceled download to %v\n", r.Path)
	case *dose.ErrorResponse:
		fmt.Printf("Received error: %s\n", r.Error)
	default:
		fmt.Printf("Received unexpected message of type %v\n", response.MessageType())
	}
}
