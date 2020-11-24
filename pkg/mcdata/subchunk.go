package mcdata

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"

	"github.com/midnightfreddie/nbt2json"
)

type SubChunk struct {
	data         []byte         // The raw sub chunk data
	Version      int            // The version of the data format (may be 1 or 8)
	StorageCount int            // Count of Block storage records (unused if version is set to 1)
	BlockStorage []BlockStorage // Zero or more concatenated Block Storage records, as specified by the count
	// (or 1 if version is set to 1).
}

type BlockStorage struct {
	version int

	blockStateIndices interface{} // The block states as indices into the palette, packed into
	// ceil(4096 / blocksPerWord) 32-bit little-endian unsigned integers.

	paletteSize uint32 // A 32-bit little-endian integer specifying the number of block states in the
	// palette.

	blockStates []BlockState // The specified number of block states in little-endian NBT format, concatenated.
}

type BlockState struct {
	TagType int
	Name    string
	Value   interface{}
}

func (b *BlockStorage) StateName(index int) (string, error) {
	state := b.blockStates[index]

	for _, t := range state.Value.([]interface{}) {
		tag := t.(map[string]interface{})
		if tag["name"] == "name" {
			return tag["value"].(string), nil
		}
	}

	return "", fmt.Errorf("reading block name: no tag found with name 'name'")
}

func NewSubChunk(data []byte) (SubChunk, error) {
	r := bytes.NewReader(data)

	version := int(readByte(r))

	switch version {
	case 1:
		log.Fatal("HANDLE SUBCHUNK TYPE 1")
		return SubChunk{}, nil
	case 8:
		// Number of BlockStorage objects to read
		storageCount := int(readBytes(r, 1)[0])

		blocks := make([]BlockStorage, storageCount)

		// Read BlockStorage data and create objects
		for i := 0; i < storageCount; i++ {
			b, err := NewBlockStorage(r)
			if err != nil {
				return SubChunk{}, fmt.Errorf("creating new block: %s", err)
			}

			blocks[i] = b
		}

		return SubChunk{
			data:         data,
			Version:      version,
			BlockStorage: blocks,
		}, nil
	default:
		panic("sub chunk had version other than 1 or 8")
	}
}

func NewBlockStorage(data *bytes.Reader) (BlockStorage, error) {
	// Version and bitsPerBlock in a single byte
	storageVersionByte := readByte(data)

	// The version (0 or 1)
	storageVersionFlag := int((storageVersionByte >> 1) & 1)

	// Number of bits used for one block state index
	bitsPerBlock := int(storageVersionByte >> 1)

	// Number of blocks per 32-bit integer
	blocksPerWord := math.Floor(float64(32 / bitsPerBlock))

	// Total count of block state indices
	indexCount := int(math.Ceil(4096/blocksPerWord)) * int(blocksPerWord)

	if 32%int(blocksPerWord) != 0 { // TODO: Handle all blocksPerword amounts https://minecraft.gamepedia.com/Bedrock_Edition_level_format
		// "For the blocksPerWord values which are not factors of 32, each 32-bit integer contains two (high) bits of padding. Block state indices are not split across words."
		// Probably need to handle: "Block state indices are *not split across words*"
		log.Fatalf("blocksPerWord is not a factor of 32")
	}

	if bitsPerBlock != 4 { // TODO: Handle all bitsPerBlock amounts https://minecraft.gamepedia.com/Bedrock_Edition_level_format
		log.Fatal("bitsPerBlock is not 4")
	}

	dataBits := NewBitReader(data)

	indices := make([]int, indexCount)
	//set := make(map[string]int) //DEBUG
	for i := 0; i < indexCount; i++ {
		// Read one block
		idxBits, err := dataBits.ReadBits(bitsPerBlock)
		if err != nil {
			return BlockStorage{}, nil
		}

		idx := int(boolsToBytes(idxBits)[0] >> 4) // TODO: see if statement above, this is specific to a bitsPerBlock value of 4. Because we are converting 4 bits to a byte, we shift it 4 bits to the right to get the correct value.
		indices[i] = idx
	}

	if dataBits.Offset() != 8 { // TODO: This does not necessarily mean things are broken
		log.Fatalf("finished reading indices part way through a byte")
	}

	// Number of blocks states in the palette
	paletteSize := binary.LittleEndian.Uint32(readBytes(data, 4))

	// Read all the remaining bytes. This is the NBT block states.
	remaining, err := ioutil.ReadAll(data)
	if err != nil {
		return BlockStorage{}, fmt.Errorf("reading remaining bytes: %s", err)
	}

	// Convert the BNT to JSON then unmarshal the JSON.
	jsn, err := nbt2json.Nbt2Json(remaining, "#")
	if err != nil {
		log.Fatal(err)
	}

	var nbtJsonData nbt2json.NbtJson
	err = json.Unmarshal(jsn, &nbtJsonData)

	if err != nil {
		log.Fatal(err)
	}

	blockStates := make([]BlockState, paletteSize)

	// Unmarshal the json states
	for i, j := range nbtJsonData.Nbt {
		var state BlockState

		err := json.Unmarshal(*j, &state)
		if err != nil {
			return BlockStorage{}, fmt.Errorf("unmarshalling block state tag: %s", err)
		}

		blockStates[i] = state
	}

	blockStorage := BlockStorage{
		version:           storageVersionFlag,
		blockStateIndices: indices,
		paletteSize:       paletteSize,
		blockStates:       blockStates,
	}

	for i := range blockStorage.blockStates {
		fmt.Println(blockStorage.StateName(i))
	}

	return blockStorage, nil
}

// func reads count byte from reader and returns, or exits the program if reader.Read() returns an error.
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
