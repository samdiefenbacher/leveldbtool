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

// subChunk is block data for a 16x16x16 area of the map.
type subChunk struct {
	data         []byte         // The raw sub chunk data
	Version      int            // The version of the data format (may be 1 or 8)
	StorageCount int            // Count of Block storage records (unused if version is set to 1)
	BlockStorage []BlockStorage // Zero or more concatenated Block Storage records, as specified by the count
	// (or 1 if version is set to 1).
}

// BlockStorage has a 'palette' (slice) of varying size containing block states and a slice of indices referring to the
// state palette.
type BlockStorage struct {
	version int

	BlockStateIndices []int // The block states as indices for statePalette.

	statePalette []Tag // Block states as key value data, each describing a block type and current state combination.
}

// Tag is JSON like data decoded from NBT blocks.
type Tag struct {
	TagType int
	Name    string
	Value   interface{}
}

func newTag(data interface{}) Tag {
	d := data.(map[string]interface{})

	return Tag{
		TagType: int(d["tagType"].(float64)),
		Name:    d["name"].(string),
		Value:   d["value"],
	}
}

// Get returns a copy of the block state from the palette using the state index at the given index.
//
// Blocks are stored column-by-column: increment Y first, Z at the end of the column and X at the end of the
// cross-section.
//
func (b *BlockStorage) Get(index int) Block {
	return b.State(b.BlockStateIndices[index])
}

// State returns the block state at the given index in the block palette. To get the state for a specific block, use
// BlockStorage.Get.
func (b *BlockStorage) State(index int) Block {
	stateData := b.statePalette[index]

	var name string

	var version int

	states := make(map[string]interface{})

	for _, t := range stateData.Value.([]interface{}) {
		tag := t.(map[string]interface{})

		switch tag["name"] {
		case "name":
			name = tag["value"].(string)
		case "version":
			version = int(tag["value"].(float64))
		case "states":
			for _, s := range tag["value"].([]interface{}) {
				st := s.(map[string]interface{})
				states[st["name"].(string)] = st["value"]
			}
		default:
			panic(fmt.Sprintf("unhandled state tag: %s", tag["name"].(string)))
		}
	}

	block := Block{
		Name:    name,
		States:  states,
		Version: version,
	}

	return block
}

func newSubChunk(data []byte) (subChunk, error) {
	r := bytes.NewReader(data)

	version := int(readByte(r))

	switch version {
	case 1:
		log.Fatal("HANDLE SUBCHUNK TYPE 1")
		return subChunk{}, nil
	case 8:
		// Number of BlockStorage objects to read
		storageCount := int(readBytes(r, 1)[0])
		blocks := make([]BlockStorage, storageCount)

		// Read BlockStorage data
		for i := 0; i < storageCount; i++ {
			b, err := readBlockStorage(r)
			if err != nil {
				return subChunk{}, fmt.Errorf("creating new block: %s", err)
			}

			blocks[i] = b
		}

		return subChunk{
			data:         data,
			Version:      version,
			BlockStorage: blocks,
		}, nil
	default:
		panic("sub chunk had version other than 1 or 8")
	}
}

func readBlockStorage(data *bytes.Reader) (BlockStorage, error) {
	// Version and bitsPerBlock in a single byte
	storageVersionByte := readByte(data)

	// The version (0 or 1)
	storageVersionFlag := int((storageVersionByte >> 1) & 1)

	// Number of bits used for one block state index
	bitsPerBlock := int(storageVersionByte >> 1)

	// Number of blocks per 32-bit integer
	blocksPerWord := math.Floor(float64(32 / bitsPerBlock))

	// Total count of block state indices
	indexCount := 4096 // int(math.Ceil(4096/blocksPerWord)) * int(blocksPerWord)

	if 32%int(blocksPerWord) != 0 { // TODO: Handle all blocksPerword amounts https://minecraft.gamepedia.com/Bedrock_Edition_level_format
		// "For the blocksPerWord values which are not factors of 32, each 32-bit integer contains two (high) bits of padding. Block state indices are not split across words."
		// Probably need to handle: "Block state indices are *not split across words*"
		// log.Fatalf("blocksPerWord value of %f is not a factor of 32", blocksPerWord)
		return BlockStorage{}, fmt.Errorf("blocksPerWord value of %f is not a factor of 32", blocksPerWord)
	}

	if bitsPerBlock != 4 { // TODO: Handle all bitsPerBlock amounts https://minecraft.gamepedia.com/Bedrock_Edition_level_format
		// log.Fatal("bitsPerBlock is not 4")
		return BlockStorage{}, fmt.Errorf("bitsPerBlock is not 4")
	}

	dataBits := NewBitReader(data)

	indices := make([]int, indexCount)
	for i := 0; i < indexCount; i++ {
		// Read one block
		idxBits, err := dataBits.ReadBits(bitsPerBlock)
		if err != nil {
			return BlockStorage{}, nil
		}

		// Index of this block's state in the palette
		idx := int(boolsToBytes(idxBits)[0] >> 4) // TODO: see if statement above, this is specific to a bitsPerBlock value of 4. Because we are converting 4 bits to a byte, we shift it 4 bits to the right to get the correct value.
		indices[i] = idx
	}

	if dataBits.Offset() != 8 { // TODO: This does not necessarily mean things are broken
		log.Fatalf("finished reading indices of size %d bits part way through a byte", bitsPerBlock)
	}

	// Number of blocks states in the palette
	paletteSize := binary.LittleEndian.Uint32(readBytes(data, 4))

	// Read all the remaining bytes. This is the NBT block states.
	remaining, err := ioutil.ReadAll(data)
	if err != nil {
		return BlockStorage{}, fmt.Errorf("reading remaining bytes: %s", err)
	}

	stateData, err := readNBTData(remaining)

	if err != nil {
		return BlockStorage{}, err
	}

	// Construct tags from empty interfaces
	blockStates := make([]Tag, paletteSize)
	for i, j := range stateData {
		blockStates[i] = newTag(j)
	}

	blockStorage := BlockStorage{
		version:           storageVersionFlag,
		BlockStateIndices: indices,
		statePalette:      blockStates,
	}

	return blockStorage, nil
}

func readNBTData(data []byte) ([]interface{}, error) {
	// Convert the BNT to JSON then unmarshal the JSON.
	jsn, err := nbt2json.Nbt2Json(data, "#")
	if err != nil {
		log.Fatal(err)
	}

	// Unmarshal JSON NBT data
	var nbtJSON struct {
		NBT []interface{} `json:"nbt"`
	}

	err = json.Unmarshal(jsn, &nbtJSON)

	if err != nil {
		return nil, err
	}

	return nbtJSON.NBT, nil
}
