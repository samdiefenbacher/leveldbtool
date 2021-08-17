package parse

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
)

const subChunkBlockCount = 4096

// StorageCount reads the initial meta data of a subchunk record and returns an integer indicating the number of block
// storage records
func StorageCount(r io.Reader) (int, error) {
	var version int8
	if err := readLittleEndian(r, &version); err != nil {
		return 0, fmt.Errorf("reading version byte: %w", err)
	}

	var storageCount int8

	switch version {
	case 1:
		storageCount = 1
	case 8:
		if err := readLittleEndian(r, &storageCount); err != nil {
			return 0, fmt.Errorf("reading storage count: %w", err)
		}
	default:
		return 0, fmt.Errorf("unhandled subchunk block storage version: '%d'", version)
	}

	return int(storageCount), nil
}

// BlockStateIndices reads a single block storage record as the integer indices into the palette. It should be called
// the number of times returned by StorageCount, after calling StorageCount.
func BlockStateIndices(r io.Reader) ([]int, error) {
	var bitsPerBlockAndVersion byte
	if err := readLittleEndian(r, &bitsPerBlockAndVersion); err != nil {
		log.Fatalf("reading version byte: %s", err)
	}

	bitsPerBlock := int(bitsPerBlockAndVersion >> 1)

	//storageVersion := int(bitsPerBlockAndVersion & 1)
	// It seems like storageVersion = 1 for the water-logged blocks storage records
	/*if storageVersion != 0 {
		return nil, fmt.Errorf("invalid block storage version %d: 0 is expected for save files", storageVersion)
	}*/

	blocksPerWord := int(math.Floor(32.0 / float64(bitsPerBlock)))
	wordCount := int(math.Ceil(subChunkBlockCount / float64(blocksPerWord)))

	indices := make([]int, subChunkBlockCount)

	i := 0

	for w := 0; w < wordCount; w++ {
		var word int32
		if err := readLittleEndian(r, &word); err != nil {
			return nil, fmt.Errorf("reading word %d from raw data: %s", w, err)
		}

		for b := 0; b < blocksPerWord && i < subChunkBlockCount; b++ {
			indices[i] = int((word >> ((i % blocksPerWord) * bitsPerBlock)) & ((1 << bitsPerBlock) - 1))
			i++
		}
	}

	return indices, nil
}

// NBT reads the remainder of a subchunk record and returns a slice of tags. It should be called after StorageCount and
// the resulting call(s) to BlockStateIndices.
func NBT(r *bytes.Reader) ([]NBTTag, error) {
	var paletteSize int32
	if err := readLittleEndian(r, &paletteSize); err != nil {
		return nil, fmt.Errorf("reading palette size bytes: %w", err)
	}

	j, err := Nbt2Json(r, int(paletteSize))
	if err != nil {
		return nil, fmt.Errorf("calling nbt2json, %w", err)
	}

	nbt := struct {
		NBT []NBTTag
	}{}
	if err := json.Unmarshal(j, &nbt); err != nil {
		return nil, fmt.Errorf("unmarshaling json, %w", err)
	}

	if len(nbt.NBT) != int(paletteSize) {
		return nil, fmt.Errorf("%d nbt records returned for palette size of %d", len(nbt.NBT), paletteSize)
	}

	return nbt.NBT, nil
}

func readLittleEndian(r io.Reader, data interface{}) error {
	return binary.Read(r, binary.ByteOrder(binary.LittleEndian), data)
}
