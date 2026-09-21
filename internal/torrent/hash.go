package torrent

import "crypto/sha256"

func HashPiece(data []byte) [32]byte {
	return sha256.Sum256(data)
}

func VerifyPiece(data []byte, expected [32]byte) bool {
	return HashPiece(data) == expected
}
