package cache

import (
	"crypto/md5"
	"encoding/hex"
)

// Hash calculates the hash of a file's content.
func Hash(content []byte) string {
	sum := md5.Sum(content)
	return hex.EncodeToString(sum[:])
}
