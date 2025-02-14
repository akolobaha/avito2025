package hasher

import (
	"crypto/sha256"
	"encoding/hex"
)

func HashString(input string) string {
	h := sha256.New()
	h.Write([]byte(input))
	hash := h.Sum(nil)
	return hex.EncodeToString(hash)
}

func ValidateHash(input, hashed string) bool {
	return HashString(input) == hashed
}
