package torrent

import (
	"errors"
	"os"
)

func ReadPiece(
	file *os.File,
	metadata FileMetadata,
	pieceIndex int,
) ([]byte, error) {
	if pieceIndex >= metadata.PieceCount {
		return []byte{}, errors.New("invalid index")
	}
	if pieceIndex < 0 {
		return []byte{}, errors.New("invalid index")
	}
	offset := int64(pieceIndex) * metadata.PieceSize
	remaining := metadata.Size - offset
	pieceSize := min(metadata.PieceSize, remaining)
	buffer := make([]byte, pieceSize)

	n, err := file.ReadAt(buffer, offset)
	if err != nil {
		return []byte{}, err
	}

	return buffer[:n], nil
}
