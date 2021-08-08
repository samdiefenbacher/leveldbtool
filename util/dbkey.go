package util

import (
	"encoding/binary"
	"encoding/hex"
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

	key := hexKey(x)
	key += hexKey(z)

	if dimension != 0 {
		key += hexKey(dimension)
	}

	key += hex.EncodeToString([]byte{47}) // 47 is the SubChunkPrefix key type tag
	key += hex.EncodeToString([]byte{byte(y)})

	return hex.DecodeString(key)
}

func hexKey(i int32) string {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, uint32(i))
	return hex.EncodeToString(b)
}
