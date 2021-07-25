package mcdata

import (
	"bytes"
	"encoding/binary"
	"fmt"
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
		panic(fmt.Sprintf("attempting to read bytes for subchunk: %s", err))
	}

	return b
}
func readByte(reader *bytes.Reader) byte {
	return readBytes(reader, 1)[0]
}

// CoordsToSubChunkKey converts x, y and z coordingates to a hexadeciaml database key
// https://github.com/midnightfreddie/McpeTool/tree/master/docs#how-to-convert-world-coordinates-to-leveldb-keys
// TODO: Int should specify dimension (e.g. end/nether). Also probably delete mcdata.go
func coordsToSubChunkKey(x, y, z, _ int) []byte {
	// The x and z origin of the chunk
	xHex := littleEndianInt(uint32(x / subChunkSize))
	zHex := littleEndianInt(uint32(z / subChunkSize))

	subChunkPrefix := littleEndianInt(subChunkPrefixDec)[:1]
	subChunkY := littleEndianInt(uint32(y / subChunkSize))[:1]

	return append(
		append(xHex, zHex...),
		append(subChunkPrefix, subChunkY...)...)
}

/*// https://minecraft.gamepedia.com/Bedrock_Edition_level_format#SubChunkPrefix_record_.281.0_and_1.2.13_formats.29
func CoordsToSubChunkByteOffsets(x, y, z int) (blockID, blockData, skyLight, blockLight int) {
	blockID = (x%16)*256 + (z%16)*16 + (y%16 + 1)

	blockMetaBase := (x%16)*128 + (z%16)*8 + (y%16)/2

	blockData = blockMetaBase + 4096 + 1
	skyLight = blockMetaBase + 6144 + 1
	blockLight = blockMetaBase + 8192 + 1

	return
}*/

func littleEndianInt(i uint32) []byte {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, i)

	return b
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
