package zeroknowledge

import (
	"crypto/sha256"
	"encoding/hex"
)

// Data represents the data structure we wish to prove ownership of.
type Data struct {
	Content string
}

// GenerateHash creates a SHA-256 hash of the data's content.
func (d *Data) GenerateHash() string {
	h := sha256.New()
	h.Write([]byte(d.Content))
	return hex.EncodeToString(h.Sum(nil))
}

// Prover sends the hash of the data to the verifier.
func Prove(d Data) string {
	return d.GenerateHash()
}

// Verifier checks if the provided hash matches the expected hash.
func Verify(providedHash, expectedHash string) bool {
	return providedHash == expectedHash
}
