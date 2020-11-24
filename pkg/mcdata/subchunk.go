package mcdata

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"math"
)

type BitReader struct {
	reader *bytes.Reader
	byte   byte
	offset byte
}

func NewBitReader(reader *bytes.Reader) BitReader {
	return BitReader{reader: reader}
}

func (r *BitReader) Offset() int {
	return int(r.offset)
}

func (r *BitReader) ReadBits(count int) ([]bool, error) {
	b := make([]bool, count)

	for i := 0; i < count; i++ {
		bit, err := r.ReadBit()
		if err != nil {
			return b, err
		}

		b[i] = bit
	}

	return b, nil
}

func (r *BitReader) ReadBit() (bool, error) {
	if r.offset == 8 {
		r.offset = 0
	}

	if r.offset == 0 {
		var err error
		if r.byte, err = r.reader.ReadByte(); err != nil {
			return false, err
		}
	}

	bit := (r.byte & (0x80 >> r.offset)) != 0

	r.offset++

	return bit, nil
}

type SubChunk struct {
	data         []byte         // The raw sub chunk data
	Version      int            // The version of the data format (may be 1 or 8)
	StorageCount int            // Count of Block storage records (unused if version is set to 1)
	BlockStorage []BlockStorage // Zero or more concatenated Block Storage records, as specified by the count
	// (or 1 if version is set to 1).
}

type BlockStorage struct {
	version      int
	bitsPerBlock int // Valid values are 1, 2, 3, 4, 5, 6, 8 and 16. Used to calculate the number of
	// blocks per 32-bit integer (blocks per word), where blocksPerWord = floor(32 / bitsPerBlock).

	blockStateIndices interface{} // The block states as indices into the palette, packed into
	// ceil(4096 / blocksPerWord) 32-bit little-endian unsigned integers.

	paletteSize uint32 // A 32-bit little-endian integer specifying the number of block states in the
	// palette.

	blockStates []interface{} // The specified number of block states in little-endian NBT format, concatenated.
}

func NewSubChunk(data []byte) (SubChunk, error) {
	r := bytes.NewReader(data)

	version := int(readByte(r))
	fmt.Println("version:", version)

	switch version {
	case 1:
		log.Fatal("HANDLE SUBCHUNK TYPE 1")
		return SubChunk{}, nil
	case 8:
		// Number of BlockStorage objects to read
		storageCount := int(readBytes(r, 1)[0])
		fmt.Println("storageCount:", storageCount)

		blocks := make([]BlockStorage, storageCount)

		// Read BlockStorage data and create objects
		for i := 0; i < storageCount; i++ {
			b, err := NewBlock(r)
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

func NewBlock(data *bytes.Reader) (BlockStorage, error) {
	storageVersionByte := readByte(data)

	storageVersionFlag := int((storageVersionByte >> 1) & 1)

	bitsPerBlock := int(storageVersionByte >> 1)

	// Number of blocks per 32-bit integer.
	blocksPerWord := math.Floor(float64(32 / bitsPerBlock))

	// Block states as indices into the palette, packed into 32-bit little-endian unsigned integers.
	indexCount := int(math.Ceil(4096/blocksPerWord)) * int(blocksPerWord)

	if 32%int(blocksPerWord) != 0 { // TODO: This probably doesn't mean things are broken:
		// "For the blocksPerWord values which are not factors of 32, each 32-bit integer contains two (high) bits of padding. Block state indices are not split across words."
		// Probably need to handle: "Block state indices are *not split across words*"
		log.Fatalf("blocksPerWord is not a factor of 32")
	}

	fmt.Println("storageVersionFlag:", storageVersionFlag)
	fmt.Println("bitsPerBlock:", bitsPerBlock)
	fmt.Println("blocksPerWord:", blocksPerWord)
	fmt.Println("indexCount:", int(math.Ceil(4096/blocksPerWord)))

	dataBits := NewBitReader(data)

	indices := make([]int, indexCount)
	//set := make(map[string]int) //DEBUG
	for i := 0; i < indexCount; i++ {
		// Read one block
		idxBits, err := dataBits.ReadBits(bitsPerBlock)
		if err != nil {
			return BlockStorage{}, nil
		}

		idx := int(boolsToBytes(idxBits)[0] >> 4)
		indices[i] = idx

		//set[fmt.Sprintf("%d", idx)] = 0 //DEBUG
	}

	if dataBits.Offset() != 8 { // TODO: This does not necessarily mean things are broken
		log.Fatalf("finished reading indices part way through a byte")
	}

	//	for k, _ := range set { //DEBUG
	//		fmt.Printf(k)
	//	}

	paletteSize := binary.LittleEndian.Uint32(readBytes(data, 4))
	fmt.Println("paletteSize:", paletteSize)

	for i := 0; i < int(paletteSize)/2; i++ {
		readByte(data)
	}

	fmt.Println("Size:", data.Size())
	fmt.Println("Len:", data.Len())

	return BlockStorage{
		version:           storageVersionFlag,
		bitsPerBlock:      bitsPerBlock,
		blockStateIndices: indices,
		paletteSize:       paletteSize,
	}, nil
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
