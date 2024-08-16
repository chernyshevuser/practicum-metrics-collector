package utils

import (
	"bytes"
	"compress/gzip"
)

func Compress(data []byte) (*bytes.Buffer, error) {
	buf := bytes.NewBuffer(nil)
	gz := gzip.NewWriter(buf)

	if _, err := gz.Write(data); err != nil {
		return nil, err
	}

	if err := gz.Close(); err != nil {
		return nil, err
	}

	return buf, nil
}
