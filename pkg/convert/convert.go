package convert

import (
	"encoding/binary"
)

const (
	subChunkPrefixDec = 47
	chunkWidth        = 16
)

// CoordsToSubChunkKey converts x, y and z coordingates to a hexadeciaml database key
// https://github.com/midnightfreddie/McpeTool/tree/master/docs#how-to-convert-world-coordinates-to-leveldb-keys
func CoordsToSubChunkKey(x, y, z, _ int) []byte {
	// The x and y origin of the chunk as a
	xHex := convertInt(uint32(x / chunkWidth))
	zHex := convertInt(uint32(z / chunkWidth))
	subChunkPrefix := convertInt(subChunkPrefixDec)[:1]
	subChunkY := convertInt(uint32(y / chunkWidth))[:1]

	return append(
		append(xHex, zHex...),
		append(subChunkPrefix, subChunkY...)...)
}

func convertInt(i uint32) []byte {
	b := make([]byte, 32)
	binary.LittleEndian.PutUint32(b, i)

	return b[:4]
}
