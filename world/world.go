package world

import (
	"fmt"
	"log"

	"github.com/danhale-git/mine/leveldb"
	"github.com/midnightfreddie/McpeTool/world"
)

const waterID = "minecraft:water"

// Worlder is implemented by mock.World and github.com/midnightfreddie/McpeTool/world.World
type Worlder interface {
	GetBlock(x, y, z, dimension int) (Block, error)
}

type World struct {
	db world.World
}

func New(path string) (*World, error) {
	w := World{}
	l, err := world.OpenWorld(path)
	if err != nil {
		log.Fatal(err)
	}

	w.db = l

	return &w, nil
}

// TODO: Don't get the sub chunk from the DB every time, cache it

// GetBlock returns the block at the given coordinates.
func (w *World) GetBlock(x, y, z, dimension int) (Block, error) {
	key, err := leveldb.SubChunkKey(
		x, y, z,
		dimension,
	)

	value, err := w.db.Get(key)
	if err != nil {
		return Block{}, fmt.Errorf("getting sub chunk with key '%s' from leveldb: %w", key, err)
	}

	sc, err := parseSubChunk(value)
	if err != nil {
		return Block{}, fmt.Errorf("decoding sub chunk value: %w", err)
	}

	voxelIndex := subChunkVoxelToIndex(x, y, z)
	blockIndex := sc.Blocks.Indices[voxelIndex]
	blockID := sc.Blocks.Palette[blockIndex].BlockID()

	waterLogged := false
	if len(sc.WaterLogged.Indices) >= voxelIndex {
		waterIndex := sc.WaterLogged.Indices[voxelIndex]
		blockID := sc.WaterLogged.Palette[waterIndex].BlockID()
		waterLogged = blockID == waterID
	}

	return Block{
		id: blockID,
		X:  x, Y: y, Z: z,
		waterLogged: waterLogged,
	}, nil
}
