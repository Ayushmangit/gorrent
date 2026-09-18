package protocol

type MessageType uint8

const (
	MessageUnknown MessageType = iota
	MessageHello
	MessageFileInfo
	MessagePieceRequest
	MessagePieceData
)

type Message struct {
	Type    MessageType
	Payload []byte
}
