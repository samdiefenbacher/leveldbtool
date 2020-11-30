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

// BlockStorage has a 'palette' (slice) of varying size containing block states and a slice of indices referring to the
// state palette.
type BlockStorage struct {
	version int

	BlockStateIndices []int // The block states as indices for statePalette.

	statePalette []Tag // Block states as key value data, each describing a block type and current state combination.
}

/*func NewBlockStorage(data *bytes.Reader) (BlockStorage, error) {

}*/

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

func readHeaderByte(data *bytes.Reader) (version, bitsPerBlock int, err error) {
	// Version and bitsPerBlock in a single byte
	storageVersionByte := readByte(data)

	// The version (0 or 1)
	version = int((storageVersionByte >> 1) & 1)

	// Number of bits used for one block state index
	bitsPerBlock, err = int(storageVersionByte>>1), nil

	return
}

func readIndices(data *bytes.Reader, bitsPerBlock int) ([]int, error) {
	// Number of blocks per 32-bit integer
	blocksPerWord := int(math.Floor(float64(32 / bitsPerBlock)))

	fmt.Println("blocksPerWord: ", blocksPerWord)

	remainderPerWord := 32 - (blocksPerWord * bitsPerBlock)

	fmt.Println("remainderPerWord: ", remainderPerWord)

	// Total count of block state indices
	indexCount := 4096

	/*//	//	//	NEW

	dataBits := NewBitReader(data)

	indices := make([]int, indexCount) // is index count correct still ? ?

	wordCount := int(math.Ceil(4096 / float64(blocksPerWord)))

	i := 0

	for w := 0; w < wordCount; w++ {
		for b := 0; b < blocksPerWord; b++ {
			// Read one block
			idxBits, err := dataBits.ReadBits(bitsPerBlock)
			if err != nil {
				return nil, nil
			}

			// Index of this block's state in the palette
			idx := int(boolsToBytes(idxBits)[0] >> 4) // TODO: see if statement above, this is specific to a bitsPerBlock value of 4. Because we are converting 4 bits to a byte, we shift it 4 bits to the right to get the correct value.
			indices[i] = idx
			i++
		}
	}*/

	//	//	//	//	OLD

	if 32%int(blocksPerWord) != 0 { // TODO: Handle all blocksPerword amounts https://minecraft.gamepedia.com/Bedrock_Edition_level_format
		// "For the blocksPerWord values which are not factors of 32, each 32-bit integer contains two (high) bits of padding. Block state indices are not split across words."
		// Probably need to handle: "Block state indices are *not split across words*"
		// log.Fatalf("blocksPerWord value of %f is not a factor of 32", blocksPerWord)
		return nil, fmt.Errorf("blocksPerWord value of %f is not a factor of 32", blocksPerWord)
	}

	if bitsPerBlock != 4 { // TODO: Handle all bitsPerBlock amounts https://minecraft.gamepedia.com/Bedrock_Edition_level_format
		// log.Fatal("bitsPerBlock is not 4")
		return nil, fmt.Errorf("bitsPerBlock is not 4")
	}

	dataBits := NewBitReader(data)

	indices := make([]int, indexCount)
	for i := 0; i < indexCount; i++ {
		// Read one block
		idxBits, err := dataBits.ReadBits(bitsPerBlock)
		if err != nil {
			return nil, nil
		}

		// Index of this block's state in the palette
		idx := int(boolsToBytes(idxBits)[0] >> 4) // TODO: see if statement above, this is specific to a bitsPerBlock value of 4. Because we are converting 4 bits to a byte, we shift it 4 bits to the right to get the correct value.
		indices[i] = idx
	}

	//	//	//

	if dataBits.Offset() != 8 { // TODO: This does not necessarily mean things are broken
		log.Fatalf("finished reading indices of size %d bits part way through a byte", bitsPerBlock)
	}

	return indices, nil
}

func readBlockStorage(data *bytes.Reader) (BlockStorage, error) {
	storageVersionFlag, bitsPerBlock, err := readHeaderByte(data)
	if err != nil {
		return BlockStorage{}, err
	}

	if storageVersionFlag != 0 {
		fmt.Printf("WARNING: block storage version is not 0: got %d - https://minecraft.gamepedia.com/Bedrock_Edition_level_format#SubChunkPrefix_record_.281.0_and_1.2.13_formats.29\n", storageVersionFlag)
	}

	fmt.Println("bitsPerBlock: ", bitsPerBlock)

	indices, err := readIndices(data, bitsPerBlock)
	if err != nil {
		return BlockStorage{}, err
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

func readIndexBlock() {

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
