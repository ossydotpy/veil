package crypto

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
)

func GenerateRandomKey() (string, error) {
	key := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return "", fmt.Errorf("could not generate random key: %w", err)
	}
	return hex.EncodeToString(key), nil
}
