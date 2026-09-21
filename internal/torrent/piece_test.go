package torrent

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPieceRead(t *testing.T) {
	tempDir := t.TempDir()
	dataBytea := []byte("ABCDEFGHIJ")

	filePath := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(filePath, dataBytea, 0o644); err != nil {
		t.Fatalf("error writing the temp file\n %v", err)
	}

	metadata := FileMetadata{
		Name:       "test.txt",
		Size:       10,
		PieceCount: 3,
		PieceSize:  4,
	}

	file, err := os.Open(filePath)
	if err != nil {
		t.Fatalf("error opening the temp file\n %v", err)
	}
	defer file.Close()

	expected := []string{
		"ABCD",
		"EFGH",
		"IJ",
	}

	for i, want := range expected {
		piece, err := ReadPiece(file, metadata, i)
		if err != nil {
			t.Fatalf("error reading piece %d", i)
		}
		if string(piece) != want {
			t.Fatalf("expected %v\ngot % v", want, string(piece))
		}
	}
}
