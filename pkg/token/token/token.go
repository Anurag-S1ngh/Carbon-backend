package token

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
)

func GenerateHash(key string) string {
	h := sha256.New()
	h.Write([]byte(key))
	hashBytes := h.Sum(nil)

	return hex.EncodeToString(hashBytes)
}

func GenerateRandomID(length uint) (string, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
}
