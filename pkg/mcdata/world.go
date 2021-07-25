package mcdata

import (
	"fmt"
	"log"
	"math"

	"github.com/midnightfreddie/McpeTool/world"
)

type World struct {
	*world.World
}

// Chunk returns the terrain data for the 16x256x16 chunk containing the given x and z block coordinates.
func (w *World) Chunk(x, z int) (*Chunk, error) {
	c := Chunk{
		X:         int(math.Floor(float64(x)/subChunkSize) * subChunkSize),
		Z:         int(math.Floor(float64(z)/subChunkSize) * subChunkSize),
		SubChunks: make([]subChunk, worldHeight/subChunkSize),
	}

	// TODO: lots of chunks are breaking
	for y := 0; y < /*len(c.SubChunks)*/ 1; y++ {
		fmt.Println("getting chunk at origin:", c.X, y*subChunkSize, c.Z, "(", x, y, z, ")")
		sc, err := w.subChunk(x, y*subChunkSize, z)

		if err != nil {
			fmt.Printf("Skipping chunk %d because error: %s\n", y, err)
			//return nil, nil
		}

		c.SubChunks[y] = sc
	}

	/*sc, err := w.subChunk(1, 40, 1)
	if err != nil {
		panic(err)
	}

	c.SubChunks[0] = sc*/

	return &c, nil
}

// SubChunk returns the terrain data for the 16x16x16 sub chunk containing the given block coordinates.
func (w *World) subChunk(x, y, z int) (subChunk, error) {
	key := coordsToSubChunkKey(x, y, z, 0)

	// Query db for sub chunk key
	b, err := w.Get(key)
	if err != nil {
		log.Println(err)
	}

	// Read data into subChunk struct
	return newSubChunk(b, x, y, z)
}
