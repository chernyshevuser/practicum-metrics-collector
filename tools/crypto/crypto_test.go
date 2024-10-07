package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCrypto(t *testing.T) {
	t.Run("positive: encode, decode, and match source with decoded", func(t *testing.T) {
		privateKey, publicKey, err := GenerateCrypto()
		require.NoError(t, err)
		require.NotEmpty(t, privateKey)
		require.NotEmpty(t, publicKey)

		sourceMessage := "some_random_phrase123@yandex.ru"
		encodedMessage, err := Encode(publicKey, sourceMessage)
		require.NoError(t, err)
		require.NotEmpty(t, encodedMessage)

		decodedMessage, err := Decode(privateKey, encodedMessage)
		require.NoError(t, err)
		require.NotEmpty(t, decodedMessage)

		assert.Equal(t, sourceMessage, decodedMessage)
	})

	t.Run("negative: private and public key mismatch", func(t *testing.T) {
		_, publicKey, err := GenerateCrypto()
		require.NoError(t, err)
		require.NotEmpty(t, publicKey)

		wrongPrivateKey, _, err := GenerateCrypto()
		require.NoError(t, err)
		require.NotEmpty(t, wrongPrivateKey)

		sourceMessage := "some_random_phrase123@yandex.ru"
		encodedMessage, err := Encode(publicKey, sourceMessage)
		require.NoError(t, err)
		require.NotEmpty(t, encodedMessage)

		decodedMessage, err := Decode(wrongPrivateKey, encodedMessage)
		assert.Error(t, err)
		assert.Empty(t, decodedMessage)
		assert.NotEqual(t, sourceMessage, decodedMessage)
	})

	t.Run("negative: empty public key", func(t *testing.T) {
		sourceMessage := "some_random_phrase123@yandex.ru"
		encodedMessage, err := Encode("", sourceMessage)
		assert.Error(t, err)
		assert.Empty(t, encodedMessage)
	})

	t.Run("negative: empty source message", func(t *testing.T) {
		_, publicKey, err := GenerateCrypto()
		require.NoError(t, err)
		require.NotEmpty(t, publicKey)

		sourceMessage := ""
		encodedMessage, err := Encode(publicKey, sourceMessage)
		assert.Error(t, err)
		assert.Empty(t, encodedMessage)
	})

	t.Run("negative: empty private key", func(t *testing.T) {
		_, publicKey, err := GenerateCrypto()
		require.NoError(t, err)
		require.NotEmpty(t, publicKey)

		sourceMessage := "some_random_phrase123@yandex.ru"
		encodedMessage, err := Encode(publicKey, sourceMessage)
		require.NoError(t, err)
		require.NotEmpty(t, encodedMessage)

		decodedMessage, err := Decode("", encodedMessage)
		assert.Error(t, err)
		assert.Empty(t, decodedMessage)
	})
}

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
