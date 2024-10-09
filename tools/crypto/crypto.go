// Package crypto provides encrypt/decrypt methods.
package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

// Encrypt encrypts a string using AES with a fixed IV.
func Encrypt(key, text string, fixedIVStr string) (string, error) {
	keyBytes := []byte(key)
	if len(keyBytes) != 16 && len(keyBytes) != 24 && len(keyBytes) != 32 {
		return "", fmt.Errorf("key length must be 16, 24, or 32 bytes")
	}

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}

	b := []byte(text)
	ciphertext := make([]byte, len(b))

	stream := cipher.NewCFBEncrypter(block, []byte(fixedIVStr))
	stream.XORKeyStream(ciphertext, b)

	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts a string using AES with a fixed IV.
func Decrypt(key, cryptoText string, fixedIVStr string) (string, error) {
	keyBytes := []byte(key)
	if len(keyBytes) != 16 && len(keyBytes) != 24 && len(keyBytes) != 32 {
		return "", fmt.Errorf("key length must be 16, 24, or 32 bytes")
	}

	ciphertext, err := base64.URLEncoding.DecodeString(cryptoText)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}

	plaintext := make([]byte, len(ciphertext))

	stream := cipher.NewCFBDecrypter(block, []byte(fixedIVStr))
	stream.XORKeyStream(plaintext, ciphertext)

	return string(plaintext), nil
}

func Sign(data []byte, key string) []byte {
	h := hmac.New(sha256.New, []byte(key))
	dst := h.Sum(data)
	return dst
}
