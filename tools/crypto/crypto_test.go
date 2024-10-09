package crypto_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/chernyshevuser/practicum-metrics-collector/tools/crypto"
	"github.com/test-go/testify/assert"
)

func Test_EncryptDecrypt(t *testing.T) {
	const (
		key        = "examplekey123456"
		fixedIVStr = "1234567890123456"
	)

	tests := []struct {
		name string
		text string
		err  error
	}{
		{
			"random phrase",
			"hello, world!",
			nil,
		},
		{
			"random hex",
			"0x70bc0dc6414eb8974bc70685f798838a87d8cce4",
			nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			encrypted, err := crypto.Encrypt(key, test.text, fixedIVStr)
			assert.Equal(t, err, test.err)
			assert.NotEqual(t, encrypted, test.text)

			decrypted, err := crypto.Decrypt(key, encrypted, fixedIVStr)
			assert.Equal(t, err, test.err)
			assert.Equal(t, decrypted, test.text)
		})
	}
}

func Test_Encrypt(t *testing.T) {
	const fixedIVStr = "1234567890123456"

	tests := []struct {
		key      string
		name     string
		text     string
		expected string
		err      error
	}{
		{
			"examplekey123456",
			"random phrase",
			"hello, world!",
			"bcsutqLSTxuL5StzzQ==",
			nil,
		},
		{
			"examplekey123456",
			"random hex",
			"0x70bc0dc6414eb8974bc70685f798838a87d8cce4",
			"NdZ16q-dXwiHoXMm2ImTOPurrms5Nu72rGds30b7Xdf1S5SGVH1Rc4ZP",
			nil,
		},
		{
			"examplekey12345",
			"random hex",
			"0x70bc0dc6414eb8974bc70685f798838a87d8cce4",
			"",
			fmt.Errorf("key length must be 16, 24, or 32 bytes"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			encrypted, err := crypto.Encrypt(test.key, test.text, fixedIVStr)
			assert.Equal(t, err, test.err)
			if test.err == nil {
				assert.Equal(t, encrypted, test.expected)
			}
		})
	}
}
func Test_Decrypt(t *testing.T) {
	const fixedIVStr = "1234567890123456"

	tests := []struct {
		key      string
		name     string
		text     string
		expected string
		err      error
	}{
		{
			"examplekey123456",
			"random phrase",
			"bcsutqLSTxuL5StzzQ==",
			"hello, world!",
			nil,
		},
		{
			"examplekey123456",
			"random hex",
			"NdZ16q-dXwiHoXMm2ImTOPurrms5Nu72rGds30b7Xdf1S5SGVH1Rc4ZP",
			"0x70bc0dc6414eb8974bc70685f798838a87d8cce4",
			nil,
		},
		{
			"",
			"random hex",
			"NdZ16q-dXwiHoXMm2ImTOPurrms5Nu72rGds30b7Xdf1S5SGVH1Rc4ZP",
			"",
			fmt.Errorf("key length must be 16, 24, or 32 bytes"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			encrypted, err := crypto.Decrypt(test.key, test.text, fixedIVStr)
			assert.Equal(t, err, test.err)
			if err == nil {
				assert.Equal(t, encrypted, test.expected)
			}
		})
	}
}

func BenchmarkEncrypt(b *testing.B) {
	const (
		key        = "thisis32bitlongpassphraseimusing"
		fixedIVStr = "1234567890123456"
		text       = "This is a test string for encryption"
	)

	for i := 0; i < b.N; i++ {
		_, err := crypto.Encrypt(key, text, fixedIVStr)
		if err != nil {
			b.Errorf("Failed to encrypt: %v", err)
		}
	}
}

func BenchmarkDecrypt(b *testing.B) {
	const (
		key        = "thisis32bitlongpassphraseimusing"
		fixedIVStr = "1234567890123456"
		text       = "This is a test string for encryption"
	)

	encryptedText, err := crypto.Encrypt(key, text, fixedIVStr)
	if err != nil {
		b.Fatalf("Failed to encrypt: %v", err)
	}

	for i := 0; i < b.N; i++ {
		_, err := crypto.Decrypt(key, encryptedText, fixedIVStr)
		if err != nil {
			b.Errorf("Failed to decrypt: %v", err)
		}
	}
}

func Test_Sign(t *testing.T) {
	data := []byte("hello world")
	key := "thisis32bitlongpassphraseimusing"

	expected := []byte{104, 101, 108, 108, 111, 32, 119, 111, 114, 108, 100, 228, 45, 53, 163, 40, 73, 82, 173, 214, 68, 18, 71, 30, 108, 79, 76, 15, 86, 176, 177, 126, 250, 222, 181, 136, 156, 2, 74, 153, 56, 207, 242}
	result := crypto.Sign(data, key)

	if !bytes.Equal(result, expected) {
		t.Errorf("Sign() = %x, want %x", result, expected)
	}
}
