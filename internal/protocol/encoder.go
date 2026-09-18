package protocol

import (
	"encoding/binary"
	"io"
)

func WriteMessage(w io.Writer, msg Message) error {
	length := uint32(1 + len(msg.Payload))

	lengthBytea := make([]byte, 4)
	binary.BigEndian.PutUint32(lengthBytea, length)

	typeBytea := []byte{byte(msg.Type)}
	n, err := w.Write(lengthBytea)
	if err != nil {
		return err
	}
	if n != len(lengthBytea) {
		return io.ErrShortWrite
	}
	n, err = w.Write(typeBytea)
	if err != nil {
		return err
	}
	if n != len(typeBytea) {
		return io.ErrShortWrite
	}
	if len(msg.Payload) > 0 {
		n, err = w.Write(msg.Payload)
		if err != nil {
			return err
		}
		if n != len(msg.Payload) {
			return io.ErrShortWrite
		}
	}
	return nil
}
