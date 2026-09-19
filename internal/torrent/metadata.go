package torrent

import (
	"errors"
	"os"
)

const DefaultPieceSize int64 = 256 * 1024

type FileMetadata struct {
	Name       string
	Size       int64
	PieceSize  int64
	PieceCount int
}

func CalculatePieceCount(fileSize, pieceSize int64) int {
	if fileSize == 0 {
		return 0
	}
	if pieceSize <= 0 {
		return 0
	}
	return int((fileSize + pieceSize - 1) / pieceSize)
}

func BuildFileMetadata(path string) (FileMetadata, error) {
	file, err := os.Open(path)
	if err != nil {
		return FileMetadata{}, err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return FileMetadata{}, err
	}
	if info.IsDir() {
		return FileMetadata{}, errors.New("is a directory")
	}

	pieceCount := CalculatePieceCount(info.Size(), DefaultPieceSize)

	metadata := FileMetadata{
		Name:       info.Name(),
		Size:       info.Size(),
		PieceSize:  DefaultPieceSize,
		PieceCount: pieceCount,
	}

	return metadata, nil
}
