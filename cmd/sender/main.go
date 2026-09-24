package main

import (
	"errors"
	"io"
	"log"
	"math"
	"net"
	"os"

	"github.com/Ayushmangit/gorrent/internal/protocol"
	"github.com/Ayushmangit/gorrent/internal/torrent"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("usage: sender <file>")
	}
	filePath := os.Args[1]

	metadata, err := torrent.BuildFileMetadata(filePath)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf(
		"sharing %s (%d bytes, %d pieces)",
		metadata.Name,
		metadata.Size,
		metadata.PieceCount,
	)

	// NOTE: validate the sizes
	if metadata.Size < 0 {
		log.Fatal(errors.New("invalid file size"))
	}

	if metadata.PieceSize <= 0 ||
		metadata.PieceSize > math.MaxUint32 {
		log.Fatal(errors.New("invalid piece size"))
	}

	if metadata.PieceCount < 0 ||
		uint64(metadata.PieceCount) > math.MaxUint32 {
		log.Fatal(errors.New("invalid piece count"))
	}

	info := protocol.FileInfo{
		Name:        metadata.Name,
		Size:        uint64(metadata.Size),
		PieceSize:   uint32(metadata.PieceSize),
		PieceCount:  uint32(metadata.PieceCount),
		PieceHashes: metadata.PieceHashes,
	}

	msg, err := protocol.NewFileInfo(info)
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	listener, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	log.Println("waiting for peer on :9000")
	conn, err := listener.Accept()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	log.Printf("peer connected: %s", conn.RemoteAddr())
	err = protocol.WriteMessage(conn, msg)
	if err != nil {
		log.Fatal(err)
	}

	for {
		requestMsg, err := protocol.ReadMessage(conn)
		if err != nil {
			if errors.Is(err, io.EOF) {
				log.Println("receiver disconnected")
				break
			}
			log.Fatal(err)
		}

		pieceIndex, err := protocol.ParsePieceRequest(requestMsg)
		if err != nil {
			log.Fatal(err)
		}

		if uint64(pieceIndex) >= uint64(metadata.PieceCount) {
			log.Fatal("invalid piece index")
		}

		piece, err := torrent.ReadPiece(
			file,
			metadata,
			int(pieceIndex),
		)
		if err != nil {
			log.Fatal(err)
		}

		response := protocol.NewPieceData(pieceIndex, piece)

		if err := protocol.WriteMessage(conn, response); err != nil {
			log.Fatal(err)
		}

		log.Printf("sent piece %d", pieceIndex)
	}
}
