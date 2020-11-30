package mcdata

import (
	"bytes"
	"fmt"
	"log"
)

// subChunk is block data for a 16x16x16 area of the map.
type subChunk struct {
	data         []byte         // The raw sub chunk data
	Version      int            // The version of the data format (may be 1 or 8)
	StorageCount int            // Count of Block storage records (unused if version is set to 1)
	BlockStorage []BlockStorage // Zero or more concatenated Block Storage records, as specified by the count
	// (or 1 if version is set to 1).
}

func newSubChunk(data []byte) (subChunk, error) {
	r := bytes.NewReader(data)

	version := int(readByte(r))

	switch version {
	case 1:
		log.Fatal("HANDLE SUBCHUNK TYPE 1")
		return subChunk{}, nil
	case 8:
		// Number of BlockStorage objects to read
		storageCount := int(readBytes(r, 1)[0])
		blocks := make([]BlockStorage, storageCount)

		// Read BlockStorage data
		for i := 0; i < storageCount; i++ {
			b, err := readBlockStorage(r)
			if err != nil {
				return subChunk{}, fmt.Errorf("creating new block: %s", err)
			}

			blocks[i] = b
		}

		return subChunk{
			data:         data,
			Version:      version,
			BlockStorage: blocks,
		}, nil
	default:
		panic("sub chunk had version other than 1 or 8")
	}
}
