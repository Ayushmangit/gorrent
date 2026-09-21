package torrent

import "testing"

func TestHash(t *testing.T) {
	data := []byte("ABCD")

	hash := HashPiece(data)

	res1 := VerifyPiece(data, hash)
	res2 := VerifyPiece([]byte("ABCE"), hash)

	if res1 != true {
		t.Fatalf("expected true got false")
	}

	if res2 != false {
		t.Fatalf("expected false got true")
	}
}
