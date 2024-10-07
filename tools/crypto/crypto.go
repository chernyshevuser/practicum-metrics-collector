package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io"
	"math/big"
	"os"
	"time"
)

// Decode decrypts a base64-encoded message using the provided private key.
func Decode(privateKeyPEM, message string) (string, error) {
	if privateKeyPEM == "" {
		return "", fmt.Errorf("private key is empty")
	}
	if message == "" {
		return "", fmt.Errorf("message is empty")
	}

	block, _ := pem.Decode([]byte(privateKeyPEM))
	if block == nil {
		return "", fmt.Errorf("failed to decode private key PEM block")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("failed to parse private key: %v", err)
	}

	decodedMessage, err := base64.StdEncoding.DecodeString(message)
	if err != nil {
		return "", fmt.Errorf("failed to decode message: %w", err)
	}

	decryptedBytes, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, decodedMessage)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt message: %w", err)
	}

	return string(decryptedBytes), nil
}

// Encode encrypts a message using the provided public key.
func Encode(publicKeyPEM, message string) (string, error) {
	if publicKeyPEM == "" {
		return "", fmt.Errorf("public key is empty")
	}
	if message == "" {
		return "", fmt.Errorf("message is empty")
	}

	block, _ := pem.Decode([]byte(publicKeyPEM))
	if block == nil {
		return "", fmt.Errorf("failed to decode public key PEM block")
	}

	certificate, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("failed to parse certificate: %v", err)
	}

	rsaPublicKey, ok := certificate.PublicKey.(*rsa.PublicKey)
	if !ok {
		return "", fmt.Errorf("unsupported public key type: %T", certificate.PublicKey)
	}

	encryptedBytes, err := rsa.EncryptPKCS1v15(rand.Reader, rsaPublicKey, []byte(message))
	if err != nil {
		return "", fmt.Errorf("failed to encrypt message: %w", err)
	}

	return base64.StdEncoding.EncodeToString(encryptedBytes), nil
}

// GenerateCrypto generates a new RSA private key and self-signed certificate.
func GenerateCrypto() (string, string, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate private key: %v", err)
	}

	var privateKeyPEM bytes.Buffer
	err = pem.Encode(&privateKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})
	if err != nil {
		return "", "", fmt.Errorf("failed to encode private key: %w", err)
	}

	certificateTemplate := &x509.Certificate{
		SerialNumber: big.NewInt(123456789),
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(10, 0, 0),
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}

	certificateBytes, err := x509.CreateCertificate(rand.Reader, certificateTemplate, certificateTemplate, &privateKey.PublicKey, privateKey)
	if err != nil {
		return "", "", fmt.Errorf("failed to create certificate: %v", err)
	}

	var certPEM bytes.Buffer
	err = pem.Encode(&certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certificateBytes,
	})
	if err != nil {
		return "", "", fmt.Errorf("failed to encode certificate: %w", err)
	}

	return privateKeyPEM.String(), certPEM.String(), nil
}

// Encrypt encrypts a string using AES.
func Encrypt(key, text string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	b := []byte(text)
	ciphertext := make([]byte, aes.BlockSize+len(b))
	iv := ciphertext[:aes.BlockSize]

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], b)

	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts a string using AES.
func Decrypt(key, cryptoText string) (string, error) {
	ciphertext, err := base64.URLEncoding.DecodeString(cryptoText)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	if len(ciphertext) < aes.BlockSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return string(ciphertext), nil
}

func Sign(data []byte, key string) []byte {
	h := hmac.New(sha256.New, []byte(key))
	dst := h.Sum(data)
	return dst
}

func LoadFromFile(path string) (string, error) {
	keyData, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read key file: %v", err)
	}
	return string(keyData), nil
}
