package protocol

import (
	"encoding/binary"
	"fmt"
	"io"
)

const MaxMessageSize = 1 << 20 // 1MiB

func ReadMessage(r io.Reader) (Message, error) {
	header := make([]byte, 4)
	_, err := io.ReadFull(r, header)
	if err != nil {
		return Message{}, fmt.Errorf("error reading header: %w", err)
	}

	length := binary.BigEndian.Uint32(header)

	if length < 1 || length > MaxMessageSize {
		return Message{}, fmt.Errorf("invalid length")
	}

	body := make([]byte, length)
	_, err = io.ReadFull(r, body)
	if err != nil {
		return Message{}, fmt.Errorf("error reading body: %w", err)
	}

	messageType := MessageType(body[0])

	msg := Message{
		Type:    messageType,
		Payload: body[1:],
	}
	return msg, nil
}
