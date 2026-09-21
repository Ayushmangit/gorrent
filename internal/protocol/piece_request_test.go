package protocol

import "testing"

func TestPieceRequest(t *testing.T) {
	original := uint32(7)
	msg := NewPieceRequest(original)

	received, err := ParsePieceRequest(msg)
	if err != nil {
		t.Fatalf("error parsing the message: %v", err)
	}

	if original != received {
		t.Fatalf("expected %v\ngot %v", original, received)
	}

	malInfo := Message{
		Type:    MessagePieceRequest,
		Payload: []byte{0, 7},
	}

	_, err = ParsePieceRequest(malInfo)

	if err == nil {
		t.Fatal("expected malformed payload to return an error")
	}

	wrongType := Message{
		Type:    MessageHello,
		Payload: []byte{0, 0, 0, 7},
	}

	_, err = ParsePieceRequest(wrongType)
	if err == nil {
		t.Fatalf("expected wrongType payload to return an error")
	}
}
