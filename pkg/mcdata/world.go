package mcdata

import (
	"encoding/binary"
	"fmt"
	"log"

	"github.com/midnightfreddie/McpeTool/world"
)

type World struct {
	*world.World
}

// SubChunk returns the terrain data for the 16x16x16 sub chunk containing the given block coordinates.
func (w *World) subChunk(x, y, z int) (subChunk, error) {
	key := coordsToSubChunkKey(x, y, z, 0)

	// Query db for sub chunk key
	b, err := w.Get(key)
	if err != nil {
		log.Println(err)
	}

	// Read data into subChunk struct
	return NewSubChunk(b)
}

// Chunk returns the terrain data for the 16x256x16 chunk containing the given x and z block coordinates.
func (w *World) Chunk(x, z int) (*Chunk, error) {
	c := Chunk{
		SubChunks: make([]subChunk, worldHeight/subChunkSize),
	}

	for y := 0; y < /*len(c.SubChunks)*/ 3; y++ { // TODO: lots of chunks are breaking
		sc, err := w.subChunk(x, y*subChunkSize, z)
		if err != nil {
			fmt.Printf("Skipping chunk %d because error: %s\n", y, err)
			//return nil, nil
		}

		c.SubChunks[y] = sc
	}

	return &c, nil
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
