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

// https://minecraft.gamepedia.com/Bedrock_Edition_level_format#SubChunkPrefix_record_.281.0_and_1.2.13_formats.29
var allowedIndexSizes = []int{1, 2, 3, 4, 5, 6, 8, 16}

func blockIndexSizeAllowed(s int) bool {
	for _, allowed := range allowedIndexSizes {
		if s == allowed {
			return true
		}
	}

	return false
}

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

func (b *BlockStorage) StateCOUNT() int {
	return len(b.statePalette)
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

func readBlockStorage(data *bytes.Reader) (BlockStorage, error) {
	fmt.Println("original data length:", data.Len())

	storageVersionFlag, bitsPerBlock, err := readHeaderByte(data)
	if err != nil {
		return BlockStorage{}, err
	}

	if storageVersionFlag != 0 {
		fmt.Printf("WARNING: block storage version is not 0: got %d - https://minecraft.gamepedia.com/Bedrock_Edition_level_format#SubChunkPrefix_record_.281.0_and_1.2.13_formats.29\n", storageVersionFlag)
	}

	fmt.Println("bitsPerBlock: ", bitsPerBlock)

	if !blockIndexSizeAllowed(bitsPerBlock) {
		return BlockStorage{}, fmt.Errorf("illegal block index bit array length of %d", bitsPerBlock)
	}

	indices, err := readIndices(data, bitsPerBlock)
	if err != nil {
		return BlockStorage{}, err
	}

	// Number of blocks states in the palette
	paletteSize := binary.LittleEndian.Uint32(readBytes(data, 4))

	fmt.Println("paletteSize:", paletteSize)

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

	//for i := 0; i < int(paletteSize); i++ {
	//	fmt.Println(i, blockStorage.State(i).Name, blockStorage.State(i).States)
	//}

	return blockStorage, nil
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
	fmt.Println("data length: ", data.Len())
	// Number of blocks per 32-bit integer
	blocksPerWord := int(math.Floor(float64(32 / bitsPerBlock)))

	fmt.Println("blocksPerWord: ", blocksPerWord)

	remainderPerWord := 32 - (blocksPerWord * bitsPerBlock)

	fmt.Println("remainderPerWord: ", remainderPerWord)

	wordCount := int(math.Ceil(4096 / float64(blocksPerWord)))
	fmt.Println("wordCount: ", wordCount)

	fmt.Println("blockCount:", blocksPerWord*wordCount)

	indices := make([]int, 4096) // is index count correct still ? ?

	unique := make(map[int]int)

	i := 0
OUTER:
	for w := 0; w < wordCount; w++ {
		// Read a 32 bit little endian int
		word := binary.LittleEndian.Uint32(readBytes(data, 4))
		wordBytes := make([]byte, 4)
		binary.BigEndian.PutUint32(wordBytes, word)

		//fmt.Printf("%b\n", wordBytes)

		wordBits := NewBitReader(bytes.NewReader(wordBytes))

		_, err := wordBits.ReadBits(remainderPerWord)
		if err != nil {
			return nil, err
		}

		for b := 0; b < blocksPerWord; b++ {
			// Read one block
			idxBits, err := wordBits.ReadBits(bitsPerBlock)
			if err != nil {
				return nil, nil
			}

			//debugString := ""

			idxBytes := boolsToBytes(idxBits)
			if len(idxBytes) < 2 {
				idxBytes = append(idxBytes, byte(0))
			}

			shift := 16 - bitsPerBlock

			idx16 := binary.BigEndian.Uint16(idxBytes) >> shift

			idx := int(idx16)

			/*for _, bit := range idxBits {
				if bit {
					debugString += fmt.Sprint(1)
				} else {
					debugString += fmt.Sprint(0)
				}
			}

			debugString += fmt.Sprintf(" - %b >> %d: %b - %b: %d", binary.BigEndian.Uint16(idxBytes), shift, idx16, idx, idx)
			fmt.Println(debugString)*/

			indices[i] = idx
			unique[idx] = 0
			i++
			if i == 4096 {
				break OUTER
			}
		}
	}

	fmt.Println("total block indices:", len(indices))
	fmt.Println("final data length:", data.Len())

	fmt.Println(unique)

	return indices, nil
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
