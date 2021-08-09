package parse

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"

	"github.com/midnightfreddie/nbt2json"
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
		log.Fatal(err)
	}

	bitsPerBlock := int(bitsPerBlockAndVersion >> 1)

	storageVersion := int(bitsPerBlockAndVersion & 1)
	if storageVersion != 0 {
		return nil, fmt.Errorf("invalid block storage version %d: 0 is expected for save files", storageVersion)
	}

	blocksPerWord := int(math.Floor(32.0 / float64(bitsPerBlock)))
	wordCount := int(math.Ceil(subChunkBlockCount / float64(blocksPerWord)))

	indices := make([]int, subChunkBlockCount)
	idx := 0

	for w := 0; w < wordCount; w++ {
		word := make([]byte, 4)
		if err := readLittleEndian(r, word); err != nil {
			log.Fatal(err)
		}

		// Might need to use a bit reader here if numbers other than 4 or 8 come up.
		switch bitsPerBlock {
		case 4:
			for _, b := range word {
				first := b >> 4
				second := (b << 4) >> 4

				indices[idx] = int(first)
				idx++
				indices[idx] = int(second)
				idx++
			}
		default:
			log.Panicf("unhandled bits per block '%d'", bitsPerBlock)
		}
	}

	return indices, nil
}

// NBT reads the remainder of a subchunk record and returns a slice of tags. It should be called after StorageCount and
// the resulting call(s) to BlockStateIndices.
func NBT(r io.Reader) ([]NBTTag, error) {
	var paletteSize int32
	if err := readLittleEndian(r, &paletteSize); err != nil {
		return nil, fmt.Errorf("reading palette size bytes: %w", err)
	}

	remainingBytes, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("reading remaining bytes, %w", err)
	}

	j, err := nbt2json.Nbt2Json(remainingBytes, "")
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
		return nil, fmt.Errorf("%d nbt records returned for palette size of %d", paletteSize, len(nbt.NBT))
	}

	return nbt.NBT, nil
}

func readLittleEndian(r io.Reader, data interface{}) error {
	return binary.Read(r, binary.ByteOrder(binary.LittleEndian), data)
}

/*for s := 0; s < storageCount; s++ {
var bitsPerBlockAndVersion byte
if err := read(r, &bitsPerBlockAndVersion); err != nil {
	log.Fatal(err)
}

bitsPerBlock := int(bitsPerBlockAndVersion >> 1)
storageVersion := int(bitsPerBlockAndVersion & 1)
fmt.Println("bitsPerBlock:", float64(bitsPerBlock)) //DEBUG
fmt.Println("storageVersion:", storageVersion)      //DEBUG

///////////////////////////////////////////////////////////////////////



//indices := make([]int, 0)

blockData := make([]byte, blockDataLength)
if err := read(r, blockData); err != nil {
	log.Fatal(err)
}

indices, x, y, z := getBlockDataIndices(r, bitsPerBlock)



fmt.Println("index count:", len(indices))

var paletteSize int32
if err := read(r, &paletteSize); err != nil {
log.Fatal(err)
}
fmt.Println("paletteSize:", paletteSize) //DEBUG

fmt.Println("reader length:", r.Len())
fmt.Println("data length:", len(data))

remainingBytes, err := io.ReadAll(r)
if err != nil {
log.Fatal(err)
}

j, err := nbt2json.Nbt2Json(remainingBytes, "")
if err != nil {
log.Fatal(err)
}

nbt := struct {
Nbt []NBTTag
}{}
//var nbtTags []NBTTag
if err := json.Unmarshal(j, &nbt); err != nil {
log.Fatal(err)
}

for i := range indices {
if i > 5 {
break
}
index := indices[i]
x, y, z := x[i], y[i], z[i]

fmt.Printf("(%d, %d, %d)", x, y, z)
fmt.Printf(" - %s\n", nbt.Nbt[index].BlockID())
}
}*/

/*testMap := make(map[int]bool)
for _, i := range indices {
	testMap[i] = true
}
for k := range testMap {
	fmt.Println(k)
}*/
