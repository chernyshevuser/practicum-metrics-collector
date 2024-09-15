package crypto

import (
	"testing"

	"github.com/test-go/testify/assert"
)

func TestEncrypt(t *testing.T) {
	key := "examplekey123456"

	tests := []struct {
		name string
		text string
	}{
		{
			"random phrase",
			"hello, world!",
		},
		{
			"random hex",
			"0x70bc0dc6414eb8974bc70685f798838a87d8cce4",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			encrypted, err := Encrypt(key, test.text)
			assert.Equal(t, err, nil)
			assert.NotEqual(t, encrypted, test.text)

			decrypted, err := Decrypt(key, encrypted)
			assert.Equal(t, err, nil)
			assert.Equal(t, decrypted, test.text)
		})
	}
}

func BenchmarkEncrypt(b *testing.B) {
	key := "thisis32bitlongpassphraseimusing"
	text := "This is a test string for encryption"

	for i := 0; i < b.N; i++ {
		_, err := Encrypt(key, text)
		if err != nil {
			b.Errorf("Failed to encrypt: %v", err)
		}
	}
}

func BenchmarkDecrypt(b *testing.B) {
	key := "thisis32bitlongpassphraseimusing"
	text := "This is a test string for encryption"

	encryptedText, err := Encrypt(key, text)
	if err != nil {
		b.Fatalf("Failed to encrypt: %v", err)
	}

	for i := 0; i < b.N; i++ {
		_, err := Decrypt(key, encryptedText)
		if err != nil {
			b.Errorf("Failed to decrypt: %v", err)
		}
	}
}
