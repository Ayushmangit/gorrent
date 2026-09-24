package main

import (
	"errors"
	"fmt"
	"log"
	"math"
	"net"
	"os"
	"path/filepath"
	"strings"

	"github.com/Ayushmangit/gorrent/internal/protocol"
	"github.com/Ayushmangit/gorrent/internal/torrent"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:9000")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	log.Printf("connected to sender: %s", conn.RemoteAddr())

	msg, err := protocol.ReadMessage(conn)
	if err != nil {
		log.Fatal(err)
	}

	info, err := protocol.ParseFileInfo(msg)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf(
		"receiving %s (%d bytes, %d pieces)",
		info.Name,
		info.Size,
		info.PieceCount,
	)

	if err := validateFileInfo(info); err != nil {
		log.Fatal(err)
	}

	if err := os.MkdirAll("downloads", 0o755); err != nil {
		log.Fatal(err)
	}

	outputPath := filepath.Join("downloads", info.Name)

	file, err := os.OpenFile(
		outputPath,
		os.O_CREATE|os.O_EXCL|os.O_RDWR,
		0o644,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	if err := file.Truncate(int64(info.Size)); err != nil {
		log.Fatal(err)
	}

	for pieceIndex := uint32(0); pieceIndex < info.PieceCount; pieceIndex++ {
		log.Printf(
			"requesting piece %d/%d",
			pieceIndex+1,
			info.PieceCount,
		)

		request := protocol.NewPieceRequest(pieceIndex)

		if err := protocol.WriteMessage(conn, request); err != nil {
			log.Fatal(err)
		}

		response, err := protocol.ReadMessage(conn)
		if err != nil {
			log.Fatal(err)
		}

		receivedIndex, data, err := protocol.ParsePieceData(response)
		if err != nil {
			log.Fatal(err)
		}

		if receivedIndex != pieceIndex {
			log.Fatalf(
				"requested piece %d, received piece %d",
				pieceIndex,
				receivedIndex,
			)
		}

		offset := uint64(pieceIndex) * uint64(info.PieceSize)

		// The final piece may be smaller than PieceSize.
		expectedSize := uint64(info.PieceSize)
		remaining := info.Size - offset

		if remaining < expectedSize {
			expectedSize = remaining
		}

		// Reject incomplete or oversized pieces.
		if uint64(len(data)) != expectedSize {
			log.Fatalf(
				"piece %d: expected %d bytes, received %d",
				pieceIndex,
				expectedSize,
				len(data),
			)
		}

		expectedHash := info.PieceHashes[pieceIndex]

		if !torrent.VerifyPiece(data, expectedHash) {
			log.Fatalf(
				"piece %d failed SHA-256 verification",
				pieceIndex,
			)
		}

		n, err := file.WriteAt(data, int64(offset))
		if err != nil {
			log.Fatal(err)
		}

		if n != len(data) {
			log.Fatal("incomplete piece write")
		}

		log.Printf(
			"downloaded piece %d/%d (%d bytes)",
			pieceIndex+1,
			info.PieceCount,
			len(data),
		)
	}

	// Flush file contents before reporting completion.
	if err := file.Sync(); err != nil {
		log.Fatal(err)
	}

	log.Printf(
		"download complete: %s (%d bytes)",
		outputPath,
		info.Size,
	)
}

func validateFileInfo(info protocol.FileInfo) error {
	// Accept only simple filenames, not paths.
	if info.Name == "" ||
		info.Name == "." ||
		info.Name == ".." ||
		filepath.Base(info.Name) != info.Name ||
		strings.ContainsAny(info.Name, `/\`) {
		return errors.New("invalid filename")
	}

	if info.Size > math.MaxInt64 {
		return errors.New("file exceeds supported size")
	}

	if info.PieceSize == 0 {
		return errors.New("piece size cannot be zero")
	}

	// Our current protocol decoder has a 4 MiB message limit.
	// Restrict piece sizes accordingly, allowing for the message
	// type byte and four-byte piece index.
	const maxPieceSize = 4*1024*1024 - 5

	if uint64(info.PieceSize) > maxPieceSize {
		return fmt.Errorf(
			"piece size %d exceeds supported limit",
			info.PieceSize,
		)
	}

	// Calculate the number of pieces using division and remainder,
	// avoiding overflow from size + pieceSize - 1.
	pieceSize := uint64(info.PieceSize)

	expectedCount := info.Size / pieceSize

	if info.Size%pieceSize != 0 {
		expectedCount++
	}

	if expectedCount != uint64(info.PieceCount) {
		return fmt.Errorf(
			"invalid piece count: expected %d, received %d",
			expectedCount,
			info.PieceCount,
		)
	}

	if uint64(len(info.PieceHashes)) != expectedCount {
		return errors.New("piece count does not match hash count")
	}

	return nil
}
