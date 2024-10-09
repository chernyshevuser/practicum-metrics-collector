package utils

import (
	"bytes"
	"compress/gzip"
	"io"
	"testing"
)

func TestCompress(t *testing.T) {
	data := []byte("This is a test string to compress. Let's see if it works.")

	compressedData, err := Compress(data)
	if err != nil {
		t.Fatalf("Compress() returned an error: %v", err)
	}

	gzReader, err := gzip.NewReader(compressedData)
	if err != nil {
		t.Fatalf("Failed to create gzip reader: %v", err)
	}
	defer gzReader.Close()

	var decompressedData bytes.Buffer
	_, err = io.Copy(&decompressedData, gzReader)
	if err != nil {
		t.Fatalf("Failed to decompress data: %v", err)
	}

	if !bytes.Equal(decompressedData.Bytes(), data) {
		t.Errorf("Decompressed data does not match original. Got: %v, want: %v", decompressedData.Bytes(), data)
	}
}
