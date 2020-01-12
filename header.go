package dose

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"io"
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

type Message interface {
	MessageType() MessageType
}

func WriteMessage(w io.Writer, m Message) (int, error) {
	encoded, err := json.Marshal(m)
	if err != nil {
		return 0, err
	}

	var messageType, length [2]byte
	binary.BigEndian.PutUint16(messageType[:], uint16(m.MessageType()))
	binary.BigEndian.PutUint16(length[:], uint16(len(encoded)))

	header := append(messageType[:], length[:]...)

	return w.Write(append(header, encoded...))
}

// Reads header from io.Reader, parses it, reads the message, and parses/returns that
// Returns a pointer to a Message type, e.g.
// switch r := response.(type) {
// case *dose.AddedResponse:
// }
func ReadMessage(r io.Reader) (Message, error) {
	var headerBytes [4]byte
	_, err := r.Read(headerBytes[:])
	if err != nil {
		return nil, err
	}

	header := ParseHeader(headerBytes)
	var m Message

	buf := make([]byte, header.Length)
	_, err = r.Read(buf)
	if err != nil {
		return nil, err
	}

	switch header.MessageType {
	case AddRequestMessage:
		m = &AddRequest{}
	case AddedResponseMessage:
		m = &AddedResponse{}
	case CancelRequestMessage:
		m = &CancelRequest{}
	case RemoveRequestMessage:
		m = &RemoveRequest{}
	case InfoRequestMessage:
		m = &InfoRequest{}
	case ServerInfoRequestMessage:
		m = &ServerInfoRequest{}
	case CanceledResponseMessage:
		m = &CanceledResponse{}
	case InfoResponseMessage:
		m = &InfoResponse{}
	case ServerInfoResponseMessage:
		m = &ServerInfoResponse{}
	case ErrorResponseMessage:
		m = &ErrorResponse{}
	case DownloadResponseMessage:
		m = &DownloadResponse{}
	case AuthRequestMessage:
		m = &AuthRequest{}
	case AuthResponseMessage:
		m = &AuthResponse{}
	default:
		return nil, errors.New("Unsupported/invalid message type")
	}

	err = json.Unmarshal(buf, &m)
	if err != nil {
		return nil, err
	}

	return m, nil
}

// I feel like generics would be really helpful here
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
