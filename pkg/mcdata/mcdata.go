package mcdata

import (
	"bytes"
	"log"
)

const (
	subChunkPrefixDec = 47
	subChunkSize      = 16
	worldHeight       = 256
)

// readBytes reads the given count of bytes from reader and returns, or exits the program if reader.Read() returns an
// error.
func readBytes(reader *bytes.Reader, count int) []byte {
	b := make([]byte, count)
	_, err := reader.Read(b)

	if err != nil {
		log.Fatalf("attempting to read bytes for subchunk: %s", err)
	}

	return b
}
func readByte(reader *bytes.Reader) byte {
	return readBytes(reader, 1)[0]
}

func boolsToBytes(t []bool) []byte {
	b := make([]byte, (len(t)+7)/8)
	for i, x := range t {
		if x {
			b[i/8] |= 0x80 >> uint(i%8)
		}
	}
	return b
}

func bytesToBools(b []byte) []bool {
	t := make([]bool, 8*len(b))
	for i, x := range b {
		for j := 0; j < 8; j++ {
			if (x<<uint(j))&0x80 == 0x80 {
				t[8*i+j] = true
			}
		}
	}
	return t
}
