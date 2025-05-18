// Package hashutil provides hash utility functions
package hashutil

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
)

// ComputeHash generates a SHA-256 hash of the given content.
func ComputeHash(content []byte) string {
	sum := sha256.Sum256(content)

	return hex.EncodeToString(sum[:])
}

// ComputeHashFromReader generates a SHA-256 hash from a reader.
func ComputeHashFromReader(reader io.Reader) (string, error) {
	hasher := sha256.New()
	if _, err := io.Copy(hasher, reader); err != nil {
		return "", err
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}
