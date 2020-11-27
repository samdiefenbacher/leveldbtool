package mcdata

type Chunk struct {
	X, Z      int
	SubChunks []subChunk
}
