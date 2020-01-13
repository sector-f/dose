package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/sector-f/dose"
)

func getServerInfo(conn net.Conn) {
	data := dose.ServerInfoRequest{}

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
	case *dose.ServerInfoResponse:
		response, _ := json.MarshalIndent(r, "", "    ")
		fmt.Println(string(response))
	case *dose.ErrorResponse:
		fmt.Printf("Received error: %s\n", r.Error)
	default:
		fmt.Printf("Received unexpected message of type %v\n", r.MessageType())
	}
}
