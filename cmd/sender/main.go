package main

import (
	"encoding/binary"
	"fmt"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", ":9000")
	if err != nil {
		panic(err)
	}

	fmt.Println("sender is connected to port ", conn.RemoteAddr())

	sendMessage(conn, "hello")
	sendMessage(conn, "I am peer 1")
	sendMessage(conn, "send me the 7th block")

	defer conn.Close()
}

func sendMessage(conn net.Conn, message string) error {
	length := uint32(len(message))
	lengthBytea := make([]byte, 4)

	binary.BigEndian.PutUint32(lengthBytea, length)
	messageBytea := []byte(message)
	_, err := conn.Write(lengthBytea)
	if err != nil {
		fmt.Println("error sending the header")
	}
	conn.Write(messageBytea)
	if err != nil {
		fmt.Println("error sending the message")
	}
	return nil
}
