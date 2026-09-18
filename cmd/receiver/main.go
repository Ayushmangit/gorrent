package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":9000")
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	fmt.Println("server is listening on port :9000")

	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		header := make([]byte, 4)
		for {
			_, err := io.ReadFull(conn, header)
			if err != nil {
				if err == io.EOF {
					fmt.Println("peer disconnected")
					break
				}
				fmt.Println("error reading the data", err)
				break
			}
			length := binary.BigEndian.Uint32(header)
			buffer := make([]byte, length)
			_, err = io.ReadFull(conn, buffer)
			if err != nil {
				if err == io.EOF {
					fmt.Println("read successfull")
					break
				}
				fmt.Println("error reading the data", err)
				break
			}
			fmt.Println(string(buffer))
		}
		conn.Close()
	}
}
