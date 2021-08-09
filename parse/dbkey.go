package parse

import (
	"encoding/binary"
	"math"
)

const (
	chunkSize = 16
)

// SubChunkKey builds the levelDB key for the sub chunk at the given x/y/z coordinates.
//
// https://minecraft.fandom.com/wiki/Bedrock_Edition_level_format#NBT_Structure
func SubChunkKey(x, z, dimension int32, y int) ([]byte, error) {
	x = int32(math.Floor(float64(x) / chunkSize))
	z = int32(math.Floor(float64(z) / chunkSize))
	y = int(math.Floor(float64(y) / chunkSize))

	key := make([]byte, 0)

	key = append(key, littleEndianBytes(x)...)
	key = append(key, littleEndianBytes(z)...)

	if dimension != 0 {
		key = append(key, littleEndianBytes(dimension)...)
	}

	key = append(key, []byte{47}...) // 47 is the SubChunkPrefix key type tag
	key = append(key, byte(y))

	return key, nil
}

func littleEndianBytes(i int32) []byte {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, uint32(i))
	return b
}
