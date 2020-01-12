package main

import (
	"fmt"
	"net"

	"github.com/sector-f/dose"
)

func download(conn net.Conn, url, filepath string) {
	data := dose.AddRequest{
		url,
		filepath,
	}

	_, err := dose.WriteMessage(conn, data)
	if err != nil {
		fmt.Println(err)
		return
	}

	response, err := dose.ReadMessage(conn)
	if err != nil {
		fmt.Println(err)
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
