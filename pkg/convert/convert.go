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
// TODO: Int should specify dimension (e.g. end/nether)
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

/*// https://minecraft.gamepedia.com/Bedrock_Edition_level_format#SubChunkPrefix_record_.281.0_and_1.2.13_formats.29
func CoordsToSubChunkByteOffsets(x, y, z int) (blockID, blockData, skyLight, blockLight int) {
	blockID = (x%16)*256 + (z%16)*16 + (y%16 + 1)

	blockMetaBase := (x%16)*128 + (z%16)*8 + (y%16)/2

	blockData = blockMetaBase + 4096 + 1
	skyLight = blockMetaBase + 6144 + 1
	blockLight = blockMetaBase + 8192 + 1

	return
}*/

func convertInt(i uint32) []byte {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, i)

	return b
}
