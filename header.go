package dose

import (
	"bytes"
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
	AuthRequestMessage        MessageType = 11
	AuthResponseMessage       MessageType = 12
)

// Convert request or response to []byte
// ready to be transmitted on wire
func MakeBody(r interface{}) ([]byte, error) {
	var messageType MessageType

	switch r.(type) {
	case AddRequest:
		messageType = AddRequestMessage
	case CancelRequest:
		messageType = CancelRequestMessage
	case RemoveRequest:
		messageType = RemoveRequestMessage
	case InfoRequest:
		messageType = InfoRequestMessage
	case ServerInfoRequest:
		messageType = ServerInfoRequestMessage
	case AddedResponse:
		messageType = AddedResponseMessage
	case CanceledResponse:
		messageType = CanceledResponseMessage
	case InfoResponse:
		messageType = InfoResponseMessage
	case ServerInfoResponse:
		messageType = ServerInfoResponseMessage
	case ErrorResponse:
		messageType = ErrorResponseMessage
	case DownloadResponse:
		messageType = DownloadResponseMessage
	case AuthRequest:
		messageType = AuthRequestMessage
	case AuthResponse:
		messageType = AuthResponseMessage
	default:
		return nil, errors.New("Invalid type")
	}

	data, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	header := DoseHeader{messageType, uint16(len(data))}
	headerBytes := header.ToBytes()
	response := append(headerBytes[:], data...)

	return response, nil
}

// I feel like generics would be really helpful here
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
	case AddedResponseMessage:
		data := AddedResponse{}

		err := json.Unmarshal(body, &data)
		if err != nil {
			return data, err
		}

		return data, nil
	case CanceledResponseMessage:
		data := CanceledResponse{}

		err := json.Unmarshal(body, &data)
		if err != nil {
			return data, err
		}

		return data, nil
	case ErrorResponseMessage:
		data := ErrorResponse{}

		err := json.Unmarshal(body, &data)
		if err != nil {
			return data, err
		}

		return data, nil
	case RemoveRequestMessage:
		data := RemoveRequest{}

		err := json.Unmarshal(body, &data)
		if err != nil {
			return data, err
		}

		return data, nil
	case InfoRequestMessage:
		data := InfoRequest{}

		err := json.Unmarshal(body, &data)
		if err != nil {
			return data, err
		}

		return data, nil
	case ServerInfoRequestMessage:
		data := ServerInfoRequest{}

		err := json.Unmarshal(body, &data)
		if err != nil {
			return data, err
		}

		return data, nil
	case InfoResponseMessage:
		data := InfoResponse{}

		err := json.Unmarshal(body, &data)
		if err != nil {
			return data, err
		}

		return data, nil
	case ServerInfoResponseMessage:
		data := ServerInfoResponse{}

		err := json.Unmarshal(body, &data)
		if err != nil {
			return data, err
		}

		return data, nil
	case DownloadResponseMessage:
		data := DownloadResponse{}

		err := json.Unmarshal(body, &data)
		if err != nil {
			return data, err
		}

		return data, nil
	case AuthRequestMessage:
		data := AuthRequest{}

		err := json.Unmarshal(body, &data)
		if err != nil {
			return data, err
		}

		return data, nil
	case AuthResponseMessage:
		data := AuthResponse{}

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

func (h DoseHeader) ToBytes() [4]byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian, h)
	bytes := buf.Bytes()

	return [4]byte{
		bytes[0],
		bytes[1],
		bytes[2],
		bytes[3],
	}
}

func ParseHeader(header [4]byte) DoseHeader {
	return DoseHeader{
		MessageType(binary.BigEndian.Uint16(header[:2])),
		binary.BigEndian.Uint16(header[2:]),
	}
}
