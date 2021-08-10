package terrain

import (
	"bytes"
	"log"

	"github.com/danhale-git/mine/parse"
)

// read block data
// iterate block data and get state - also calculate position

// SubChunk is the parsed data for one 16x16 subchunk. A palette including all block states in the subchunk is indexed
// by a slice of integers (one for each block) to determine the state and block id for each block in the palette.
type SubChunk struct {
	StatePalette []parse.NBTTag // The palette of block states indexed by BlockStorage to get block details
	BlockStorage []int          // Every block in the 16x16 chunk extending first on the Y axis then Z followed by X
	//waterStorage []int    // Whether blocks are waterlogged, mapped in the same way as BlockStorage
}

func NewSubChunk(data []byte) (*SubChunk, error) {
	r := bytes.NewReader(data)
	s := SubChunk{}
	c, err := parse.StorageCount(r)
	if err != nil {
		return nil, err
	}

	// https://minecraft.fandom.com/wiki/Bedrock_Edition_level_format
	// In the majority of cases, there is only one storage record.
	// A second record may be present to indicate block water-logging.
	switch c {
	case 0:
		panic("block storage count is 0")
	case 1:
		s.BlockStorage, err = parse.BlockStateIndices(r)
		if err != nil {
			return nil, err
		}
	case 2:
		panic("HANDLE WATER LOGGED BLOCKS")
		// consider converting to a slice of bools indicating if a water block was indexed
		// See SubChunk.waterStorage
	default:
		log.Panicf("unhandled storage count: %d", c)
	}

	s.StatePalette, err = parse.NBT(r)
	if err != nil {
		return nil, err
	}

	return &s, nil
}
