package main

import (
	"fmt"
	"net"
	"os"

	"github.com/sector-f/dose"
)

func download(url, filepath string) {
	data := dose.AddRequest{
		url,
		filepath,
	}

	conn, err := net.Dial("unix", "/tmp/dose.socket")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer conn.Close()

	_, err = dose.WriteMessage(conn, data)
	if err != nil {
		fmt.Println("Error writing")
		fmt.Println(err)
		return
	}

	response, err := dose.ReadMessage(conn)
	if err != nil {
		fmt.Println("Error reading")
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
