package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
)

type Engine struct {
	key []byte
}

func NewEngine(keyHex string) (*Engine, error) {
	key, err := hex.DecodeString(keyHex)
	if err != nil {
		return nil, fmt.Errorf("invalid master key format (must be hex): %w", err)
	}

	if len(key) != 32 {
		return nil, fmt.Errorf("invalid key length: expected 32 bytes, got %d", len(key))
	}
	return &Engine{key: key}, nil
}

func (e *Engine) Encrypt(value string) (string, error) {
	plaintext := []byte(value)

	block, err := aes.NewCipher(e.key)
	if err != nil {
		return "", fmt.Errorf("error creating aes block cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("error setting gcm mode: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("error generating the nonce: %w", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return hex.EncodeToString(ciphertext), nil
}

func (e *Engine) Decrypt(encryptedValue string) (string, error) {
	ciphertext, err := hex.DecodeString(encryptedValue)
	if err != nil {
		return "", fmt.Errorf("error decoding hex: %w", err)
	}

	block, err := aes.NewCipher(e.key)
	if err != nil {
		return "", fmt.Errorf("error creating aes block cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("error setting gcm mode: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("error decrypting: %w", err)
	}

	return string(plaintext), nil
}
