package dose

import (
	"encoding/binary"
	"encoding/json"
	"errors"
)

type MessageType uint16

const (
	AddRequestMessage         MessageType = 0
	CancelRequestMessage      MessageType = 1
	RemoveRequestMessage      MessageType = 2
	InfoRequestMessage        MessageType = 3
	ServerInfoRequestMessage  MessageType = 4
	AddedResponseMessage      MessageType = 5
	CanceledResponseMessage   MessageType = 6
	InfoResponseMessage       MessageType = 7
	ServerInfoResponseMessage MessageType = 8
	ErrorResponseMessage      MessageType = 9
	DownloadResponseMessage   MessageType = 10
)

func ParseBody(messageType MessageType, body []byte) (interface{}, error) {
	switch messageType {
	case AddRequestMessage:
		data := AddRequest{}

		err := json.Unmarshal(body, &data)
		if err != nil {
			return data, err
		}

		return data, nil
	case CancelRequestMessage:
		data := CancelRequest{}

		err := json.Unmarshal(body, &data)
		if err != nil {
			return data, err
		}

		return data, nil
	default:
		return nil, errors.New("Invalid message type")
	}
}

type DoseHeader struct {
	MessageType MessageType
	Length      uint16
}

func ParseHeader(header [4]byte) DoseHeader {
	return DoseHeader{
		MessageType(binary.BigEndian.Uint16(header[:2])),
		binary.BigEndian.Uint16(header[2:]),
	}
}
