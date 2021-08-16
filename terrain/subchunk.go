package terrain

import (
	"bytes"
	"fmt"
	"log"

	"github.com/danhale-git/mine/parse"
)

// read block data
// iterate block data and get state - also calculate position

// SubChunk is the parsed data for one 16x16 subchunk. A palette including all block states in the subchunk is indexed
// by a slice of integers (one for each block) to determine the state and block id for each block in the palette.
type SubChunk struct {
	Blocks      BlockStorage
	WaterLogged BlockStorage

	//StatePalette []parse.NBTTag // The palette of block states indexed by BlockStorage to get block details
	//BlockStorage []int          // Every block in the 16x16 chunk extending first on the Y axis then Z followed by X
	////waterStorage []int    // Whether blocks are waterlogged, mapped in the same way as BlockStorage
}

type BlockStorage struct {
	Indices []int          // An index into the palette for each block in the sub chunk
	Palette []parse.NBTTag // A palette of block types and states
}

func NewSubChunk(data []byte) (*SubChunk, error) {
	r := bytes.NewReader(data)
	s := SubChunk{}
	c, err := parse.StorageCount(r)
	if err != nil {
		return nil, err
	}

	fmt.Println("storage count:", c)

	fmt.Println("after reading storage count:", r.Len())

	// https://minecraft.fandom.com/wiki/Bedrock_Edition_level_format
	// In the majority of cases, there is only one storage record.
	// A second record may be present to indicate block water-logging.
	switch c {
	case 0:
		panic("block storage count is 0")
	case 1:
		s.Blocks.Indices, err = parse.BlockStateIndices(r)
		if err != nil {
			return nil, fmt.Errorf("parsing block storage indices: %s", err)
		}

		s.Blocks.Palette, err = parse.NBT(r)
		if err != nil {
			return nil, fmt.Errorf("parsing nbt data: %s", err)
		}
	case 2:
		s.Blocks.Indices, err = parse.BlockStateIndices(r)
		if err != nil {
			return nil, fmt.Errorf("parsing block storage indices: %s", err)
		}

		s.Blocks.Palette, err = parse.NBT(r)
		if err != nil {
			return nil, fmt.Errorf("parsing nbt data: %s", err)
		}

		s.WaterLogged.Indices, err = parse.BlockStateIndices(r)
		if err != nil {
			return nil, fmt.Errorf("parsing water logged indices: %s", err)
		}

		s.WaterLogged.Palette, err = parse.NBT(r)
		if err != nil {
			return nil, fmt.Errorf("parsing nbt data: %s", err)
		}

		fmt.Println("after water logged BNT:", r.Len())
	default:
		log.Panicf("unhandled storage count: %d", c)
	}

	return &s, nil
}
