package main

import (
	"encoding/binary"
	"encoding/json"
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

	encoded, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var messageType [2]byte
	binary.BigEndian.PutUint16(messageType[:], uint16(dose.AddRequestMessage))

	var length [2]byte
	binary.BigEndian.PutUint16(length[:], uint16(len(encoded)))

	header := append(messageType[:], length[:]...)

	conn, err := net.Dial("unix", "/tmp/dose.socket")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer conn.Close()

	_, err = conn.Write(header)
	if err != nil {
		fmt.Println(err)
	}

	_, err = conn.Write(encoded)
	if err != nil {
		fmt.Println(err)
	}
}
