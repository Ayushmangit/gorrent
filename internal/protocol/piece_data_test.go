package protocol

import (
	"bytes"
	"testing"
)

func TestPieceData(t *testing.T) {
	index := uint32(7)
	dataBytes := []byte("hello gorrent")

	msg := NewPieceData(index, dataBytes)
	if msg.Type != MessagePieceData {
		t.Fatalf(
			"expected message type %v, got %v",
			MessagePieceData,
			msg.Type,
		)
	}

	parsedIndex, data, err := ParsePieceData(msg)
	if err != nil {
		t.Fatalf("error parsing the msg\n %v", err)
	}

	if index != parsedIndex {
		t.Fatalf("expected %v\ngot%v", index, parsedIndex)
	}
	if !bytes.Equal(dataBytes, data) {
		t.Fatalf("expected %s\ngot%s", dataBytes, data)
	}
}
