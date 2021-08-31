package world

import (
	"fmt"
	"log"
	"math"

	"github.com/danhale-git/mine/leveldb"
	"github.com/midnightfreddie/McpeTool/world"
)

const waterID = "minecraft:water"

// BlockAPI modifies block data.
type BlockAPI interface {
	GetBlock(x, y, z, dimension int) (Block, error)
}

// LevelDB returns data from a leveldb database.
type LevelDB interface {
	Get(key []byte) ([]byte, error)
}

type World struct {
	db        LevelDB
	subChunks map[struct{ x, y, z, d int }]*subChunkData
}

func New(path string) (*World, error) {
	w := World{}
	w.subChunks = make(map[struct{ x, y, z, d int }]*subChunkData)
	l, err := world.OpenWorld(path)
	if err != nil {
		log.Fatal(err)
	}

	w.db = &l

	return &w, nil
}

// TODO: Don't get the sub chunk from the DB every time, cache it

// GetBlock returns the block at the given coordinates.
func (w *World) GetBlock(x, y, z, dimension int) (Block, error) {
	origin := subChunkOrigin(x, y, z, dimension)

	var sc *subChunkData
	var ok bool

	if sc, ok = w.subChunks[origin]; !ok {
		key, err := leveldb.SubChunkKey(
			x, y, z,
			dimension,
		)

		value, err := w.db.Get(key)
		if err != nil {

			// TODO: Make a PR to give this error a type - https://github.com/midnightfreddie/goleveldb/blob/fb12d34a9c1f2c7615bb9b258d09400cd315502f/leveldb/errors/errors.go#L19

			if err.Error() == "leveldb: not found" {
				return Block{}, &SubChunkNotSavedError{origin}
			}
			return Block{}, fmt.Errorf("getting sub chunk with key '%x': %w", key, err)
		}

		sc, err = parseSubChunk(value)
		if err != nil {
			return Block{}, fmt.Errorf("decoding sub chunk value: %w", err)
		}

		w.subChunks[origin] = sc
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
		ID: blockID,
		X:  x, Y: y, Z: z,
		waterLogged: waterLogged,
	}, nil
}

func subChunkOrigin(x, y, z, d int) struct{ x, y, z, d int } {
	return struct{ x, y, z, d int }{
		int(math.Floor(float64(x) / 16)),
		int(math.Floor(float64(y) / 16)),
		int(math.Floor(float64(z) / 16)),
		d,
	}
}

// SubChunkNotSavedError is returned if a requested sub chunk is not present in the world database.
type SubChunkNotSavedError struct {
	origin struct{ x, y, z, d int }
}

// TODO: State the dimension in the error message, when dimensions are supported

func (e *SubChunkNotSavedError) Error() string {
	return fmt.Sprintf("chunk with origin %d %d %d is not stored in this world database",
		e.origin.x, e.origin.y, e.origin.z)
}

// Is implements Is(error) to support errors.Is()
func (e *SubChunkNotSavedError) Is(tgt error) bool {
	_, ok := tgt.(*SubChunkNotSavedError)
	return ok
}
