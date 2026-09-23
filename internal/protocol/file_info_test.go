package protocol

import (
	"bytes"
	"encoding/binary"
	"reflect"
	"testing"
)

func TestNewFileInfo(t *testing.T) {
	var hash [32]byte
	hash[0] = 0xAA
	hash[1] = 0xBB

	info := FileInfo{
		Name:        "a.txt",
		Size:        100,
		PieceSize:   256,
		PieceCount:  1,
		PieceHashes: [][32]byte{hash},
	}

	msg, err := NewFileInfo(info)
	if err != nil {
		t.Fatalf("NewFileInfo returned error: %v", err)
	}
	if msg.Type != MessageFileInfo {
		t.Fatalf("expected FileInfo message,got %v", msg.Type)
	}

	if len(msg.Payload) != 55 {
		t.Fatalf("expected payload length 55,got %d", len(msg.Payload))
	}

	nameLength := binary.BigEndian.Uint16(msg.Payload[0:2])
	if nameLength != 5 {
		t.Fatalf("expected name length 5, got %d", nameLength)
	}

	name := string(msg.Payload[2:7])
	if name != "a.txt" {
		t.Fatalf("expected a.txt, got %s", name)
	}

	size := binary.BigEndian.Uint64(msg.Payload[7:15])
	if size != 100 {
		t.Fatalf("expected 100, got %d", size)
	}

	pieceSize := binary.BigEndian.Uint32(msg.Payload[15:19])
	if pieceSize != 256 {
		t.Fatalf("expected 256, got %d", pieceSize)
	}
	pieceCount := binary.BigEndian.Uint32(msg.Payload[19:23])
	if pieceCount != 1 {
		t.Fatalf("expected 1, got %d", pieceCount)
	}
	receivedHash := msg.Payload[23:55]
	if !bytes.Equal(receivedHash, hash[:]) {
		t.Fatalf("expected %v,got %v", hash, receivedHash)
	}
}

func TestFileInfoRoundTrip(t *testing.T) {
	var hash1 [32]byte
	hash1[0] = 0xAA
	hash1[1] = 0xBB

	var hash2 [32]byte
	hash2[0] = 0xCC
	hash2[1] = 0xDD

	original := FileInfo{
		Name:        "movie.mp4",
		Size:        500000,
		PieceSize:   262144,
		PieceCount:  2,
		PieceHashes: [][32]byte{hash1, hash2},
	}

	msg, err := NewFileInfo(original)
	if err != nil {
		t.Fatalf("NewFileInfo: %v", err)
	}
	parsed, err := ParseFileInfo(msg)
	if err != nil {
		t.Fatalf("NewFileInfo: %v", err)
	}

	if original.Name != parsed.Name {
		t.Fatalf("expected %s,got %s", original.Name, parsed.Name)
	}

	if original.Size != parsed.Size {
		t.Fatalf("expected %d,got %d", original.Size, parsed.Size)
	}

	if uint64(original.PieceCount) != uint64(parsed.PieceCount) {
		t.Fatalf("expected %d,got %d", original.PieceCount, parsed.PieceCount)
	}

	if uint64(original.PieceSize) != uint64(parsed.PieceSize) {
		t.Fatalf("expected %d,got %d", original.PieceSize, parsed.PieceSize)
	}

	if !reflect.DeepEqual(original.PieceHashes, parsed.PieceHashes) {
		t.Fatalf(
			"expected hashes %v, got %v",
			original.PieceHashes,
			parsed.PieceHashes,
		)
	}
}

func TestParseFileInfoMalformedPieceCount(t *testing.T) {
	var hash [32]byte
	hash[0] = 0xAA
	hash[1] = 0xBB

	info := FileInfo{
		Name:        "a.txt",
		Size:        100,
		PieceSize:   256,
		PieceCount:  1,
		PieceHashes: [][32]byte{hash},
	}

	msg, err := NewFileInfo(info)
	if err != nil {
		t.Fatalf("NewFileInfo returned error: %v", err)
	}

	// Layout for "a.txt":
	//
	// [name len][name ][size    ][piece size][piece count][hash]
	//  0:2       2:7   7:15      15:19       19:23       23:55
	//
	// The message really contains only 1 hash,
	// but we'll lie and claim that there are 2 pieces.
	binary.BigEndian.PutUint32(msg.Payload[19:23], 2)

	_, err = ParseFileInfo(msg)

	if err == nil {
		t.Fatal("expected malformed FILE_INFO to return an error")
	}
}
