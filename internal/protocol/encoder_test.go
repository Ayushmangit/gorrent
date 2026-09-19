package protocol

import (
	"bytes"
	"encoding/binary"
	"testing"
)

func TestWriteMessageHello(t *testing.T) {
	var buf bytes.Buffer

	msg := Message{
		Type: MessageHello,
	}

	err := WriteMessage(&buf, msg)
	if err != nil {
		t.Fatalf("write message returned error : %v\n", err)
	}

	expected := []byte{0, 0, 0, 1, 1}

	if !bytes.Equal(buf.Bytes(), expected) {
		t.Fatalf("Writing hello failed\nexpected: %v\n Got: %v\n", expected, buf.Bytes())
	}
}

func TestWriteMessagePieceRequest(t *testing.T) {
	var buf bytes.Buffer

	payload := uint32(7)
	payloadBytea := make([]byte, 4)
	binary.BigEndian.PutUint32(payloadBytea, payload)

	msg := Message{
		Type:    MessagePieceRequest,
		Payload: payloadBytea,
	}

	if err := WriteMessage(&buf, msg); err != nil {
		t.Fatalf("failed writing the piece request message\n %v", err)
	}

	expected := []byte{
		0, 0, 0, 5,
		3,
		0, 0, 0, 7,
	}

	if !bytes.Equal(expected, buf.Bytes()) {
		t.Fatalf(
			"writing request message failed\nexpected: %v\ngot: %v",
			expected,
			buf.Bytes(),
		)
	}
}
