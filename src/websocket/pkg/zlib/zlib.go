package zlib

import (
	zlib "github.com/4kills/go-libdeflate/v2"
)

func Compress(data []byte, level int) ([]byte, error) {
	c, err := zlib.NewCompressorLevel(2)
	if err != nil {
		return nil, err
	}
	compressed := make([]byte, len(data))

	n, _, err := c.Compress(data, compressed, zlib.ModeZlib)
	if err != nil {
		return nil, err
	}
	compressed = compressed[:n]
	defer c.Close()

	return compressed, nil
}
