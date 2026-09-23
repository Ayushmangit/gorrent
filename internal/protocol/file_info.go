package protocol

import (
	"encoding/binary"
	"errors"
)

type FileInfo struct {
	Name        string
	Size        uint64
	PieceSize   uint32
	PieceCount  uint32
	PieceHashes [][32]byte
}

func NewFileInfo(info FileInfo) (Message, error) {
	if len(info.Name) > 65535 {
		return Message{}, errors.New("filename too long")
	}
	if int(info.PieceCount) != len(info.PieceHashes) {
		return Message{}, errors.New("hashes and piece count dont match")
	}

	payloadSize := 18 + len(info.Name) + len(info.PieceHashes)*32
	payload := make([]byte, payloadSize)

	offset := 0

	binary.BigEndian.PutUint16(
		payload[offset:offset+2],
		uint16(len(info.Name)),
	)
	offset += 2
	copy(
		payload[offset:offset+len(info.Name)],
		[]byte(info.Name),
	)
	offset += len(info.Name)

	binary.BigEndian.PutUint64(
		payload[offset:offset+8],
		info.Size,
	)
	offset += 8

	binary.BigEndian.PutUint32(
		payload[offset:offset+4],
		info.PieceSize,
	)
	offset += 4

	binary.BigEndian.PutUint32(
		payload[offset:offset+4],
		info.PieceCount,
	)
	offset += 4

	for _, hash := range info.PieceHashes {
		copy(payload[offset:offset+32], hash[:])
		offset += 32
	}

	return Message{
		Type:    MessageFileInfo,
		Payload: payload,
	}, nil
}

func ParseFileInfo(msg Message) (FileInfo, error) {
	if msg.Type != MessageFileInfo {
		return FileInfo{}, errors.New("invalid message type")
	}
	if len(msg.Payload) < 2 {
		return FileInfo{}, errors.New("payload too short for filename length")
	}
	offset := 0

	nameLength := int(binary.BigEndian.Uint16(
		msg.Payload[offset : offset+2],
	))
	offset += 2

	if len(msg.Payload) < offset+nameLength {
		return FileInfo{}, errors.New("payload too short for filename")
	}

	name := string(
		msg.Payload[offset : offset+nameLength],
	)
	offset += nameLength

	if len(msg.Payload) < offset+16 {
		return FileInfo{}, errors.New("payload too short for file metadata")
	}

	size := binary.BigEndian.Uint64(
		msg.Payload[offset : offset+8],
	)
	offset += 8

	pieceSize := binary.BigEndian.Uint32(
		msg.Payload[offset : offset+4],
	)
	offset += 4

	pieceCount := binary.BigEndian.Uint32(
		msg.Payload[offset : offset+4],
	)
	offset += 4

	expectedHashBytes := uint64(pieceCount) * 32
	actualHashBytes := uint64(len(msg.Payload) - offset)

	if actualHashBytes != expectedHashBytes {
		return FileInfo{}, errors.New("invalid piece hash data length")
	}
	pieceHashes := make([][32]byte, 0, pieceCount)

	for range pieceCount {
		var hash [32]byte
		copy(hash[:], msg.Payload[offset:offset+32])
		offset += 32
		pieceHashes = append(pieceHashes, hash)
	}

	if pieceSize == 0 {
		return FileInfo{}, errors.New("piece size cannot be zero")
	}

	return FileInfo{
		Name:        name,
		Size:        size,
		PieceSize:   pieceSize,
		PieceCount:  pieceCount,
		PieceHashes: pieceHashes,
	}, nil
}
