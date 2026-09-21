package protocol

import (
	"bytes"
	"encoding/binary"
	"testing"
)

func TestReadMessageHello(t *testing.T) {
	var buf bytes.Buffer

	original := Message{
		Type: MessageHello,
	}

	if err := WriteMessage(&buf, original); err != nil {
		t.Fatal("error writing hello\n")
	}
	received, err := ReadMessage(&buf)
	if err != nil {
		t.Fatalf("error reading message\n %v", err)
	}

	if received.Type != MessageHello {
		t.Fatalf(
			"expected type %v, got %v",
			MessageHello,
			received.Type,
		)
	}

	if !bytes.Equal(received.Payload, original.Payload) {
		t.Fatalf(
			"expected payload %v, got %v",
			original.Payload,
			received.Payload,
		)
	}
}

func TestReadPieceRequestMessage(t *testing.T) {
	var buf bytes.Buffer

	payload := make([]byte, 4)
	binary.BigEndian.PutUint32(payload, 4)

	original := Message{
		Type:    MessagePieceRequest,
		Payload: payload,
	}

	if err := WriteMessage(&buf, original); err != nil {
		t.Fatalf("error writing the message: %v", err)
	}

	received, err := ReadMessage(&buf)
	if err != nil {
		t.Fatalf("error reading the message: %v", err)
	}

	if received.Type != MessagePieceRequest {
		t.Fatalf(
			"expected type %v, got %v",
			MessagePieceRequest,
			received.Type,
		)
	}

	if !bytes.Equal(received.Payload, original.Payload) {
		t.Fatalf(
			"expected payload %v, got %v",
			original.Payload,
			received.Payload,
		)
	}
}
