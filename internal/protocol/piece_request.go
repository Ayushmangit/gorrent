package protocol

import (
	"encoding/binary"
	"errors"
)

func NewPieceRequest(pieceIndex uint32) Message {
	payload := make([]byte, 4)
	binary.BigEndian.PutUint32(payload, pieceIndex)

	return Message{
		Type:    MessagePieceRequest,
		Payload: payload,
	}
}

func ParsePieceRequest(msg Message) (uint32, error) {
	if msg.Type != MessagePieceRequest {
		return 0, errors.New("invalid type")
	}

	if len(msg.Payload) != 4 {
		return 0, errors.New("invalid payload need exactly 4 bytes")
	}
	requestedPieceIndex := binary.BigEndian.Uint32(msg.Payload)
	return requestedPieceIndex, nil
}
