package protocol

import (
	"encoding/binary"
	"errors"
)

func NewPieceData(pieceIndex uint32, data []byte) Message {
	PieceData := make([]byte, 4+len(data))
	binary.BigEndian.PutUint32(PieceData, pieceIndex)
	copy(PieceData[4:], data)
	return Message{
		Type:    MessagePieceData,
		Payload: PieceData,
	}
}

func ParsePieceData(msg Message) (uint32, []byte, error) {
	if msg.Type != MessagePieceData {
		return 0, nil, errors.New("invalid message type")
	}
	if len(msg.Payload) < 5 {
		return 0, []byte{}, errors.New("invalid payload length")
	}
	indexBytes := msg.Payload[:4]
	index := binary.BigEndian.Uint32(indexBytes)
	data := msg.Payload[4:]

	return index, data, nil
}
