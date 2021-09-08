package subchunk

import (
	"log"

	"github.com/danhale-git/mine/nbt"
)

// Data is the parsed data for one 16x16 subchunk. A palette including all block states in the subchunk is indexed
// by a slice of integers (one for each block) to determine the state and block id for each block in the palette.
type Data struct {
	Blocks      blockStorage
	WaterLogged blockStorage
}

type blockStorage struct {
	Indices []int        // An index into the palette for each block in the sub chunk
	Palette []nbt.NBTTag // A palette of block types and states
}

func (d *Data) BlockState(x, y, z int) (nbt.NBTTag, bool) {
	voxelIndex := voxelToIndex(x, y, z)

	waterLogged := false
	if len(d.WaterLogged.Indices) > 0 && len(d.WaterLogged.Indices) >= voxelIndex {
		waterIndex := d.WaterLogged.Indices[voxelIndex]
		blockID := d.WaterLogged.Palette[waterIndex].BlockID()
		waterLogged = blockID == WaterID
	}

	blockIndex := d.Blocks.Indices[voxelIndex]

	return d.Blocks.Palette[blockIndex], waterLogged
}

// voxelToIndex returns the block storage index from the given sub chunk x y and z coordinates.
func voxelToIndex(x, y, z int) int {
	if x > 15 || y > 15 || z > 15 {
		log.Panicf("coordinates %d %d %d are invalid: sub chunk cooridnates may not exceed 0-15", x, y, z)
	}
	return y + z*16 + x*16*16
}

// indexToVoxel returns the world x y z offset from the sub chunk root for the given block storage index.
func indexToVoxel(i int) (x, y, z int) {
	x = (i >> 8) & 15
	y = i & 15
	z = (i >> 4) & 15

	return
}
