package util

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"math"
)

func read(r io.Reader, data interface{}) error {
	return binary.Read(r, binary.ByteOrder(binary.LittleEndian), data)
}

type NBTTag struct {
	Type  byte
	Name  string
	Value interface{}
}

func ParseSubChunk(data []byte) {
	r := bytes.NewReader(data)
	fmt.Println("reader length:", r.Len())

	var version int8
	if err := read(r, &version); err != nil {
		log.Fatal(err)
	}
	fmt.Println("version:", version) //DEBUG

	var storageCount int8 = 1

	switch version {
	case 1:
	case 8:
		if err := read(r, &storageCount); err != nil {
			log.Fatal(err)
		}
	default:
		log.Panicf("Unhandled storage version: '%d'", version)
	}
	fmt.Println("storageCount:", storageCount) //DEBUG

	var bitsPerBlock, storageVersion, blocksPerWord int

	var i int8
	for i = 0; i < storageCount; i++ {
		var bitsPerBlockAndVersion byte
		if err := read(r, &bitsPerBlockAndVersion); err != nil {
			log.Fatal(err)
		}

		bitsPerBlock = int(bitsPerBlockAndVersion >> 1)
		storageVersion = int(bitsPerBlockAndVersion & 1)
		fmt.Println("bitsPerBlock:", float64(bitsPerBlock)) //DEBUG
		fmt.Println("storageVersion:", storageVersion)      //DEBUG

		blocksPerWord = int(math.Floor(32.0 / float64(bitsPerBlock)))

		fmt.Println("blocksPerWord:", blocksPerWord) //DEBUG
	}

	wordCount := int(math.Ceil(4096 / float64(blocksPerWord)))
	fmt.Println("wordCount:", wordCount) //DEBUG

	indices := make([]int, 0)

	for w := 0; w < wordCount; w++ {
		word := make([]byte, 4)
		if err := read(r, word); err != nil {
			log.Fatal(err)
		}

		indices = append(indices, getBlockDataIndices(word, bitsPerBlock)...)
	}

	/*testMap := make(map[int]bool)
	for _, i := range indices {
		testMap[i] = true
	}
	for k := range testMap {
		fmt.Println(k)
	}*/

	fmt.Println("index count:", len(indices))

	var paletteSize int32
	if err := read(r, &paletteSize); err != nil {
		log.Fatal(err)
	}
	fmt.Println("paletteSize:", paletteSize) //DEBUG

	fmt.Println("reader length:", r.Len())
	fmt.Println("data length:", len(data))
}

func getBlockDataIndices(word []byte, bitsPerBlock int) []int {
	indices := make([]int, 0)

	// Might need to use a bit reader here if numbers other than 4 or 8 come up
	switch bitsPerBlock {
	case 4:
		for _, b := range word {
			first := b >> 4
			second := (b << 4) >> 4
			indices = append(indices, int(first), int(second))
		}
	default:
		log.Panicf("unhandled bits per block '%d'", bitsPerBlock)
	}

	return indices
}

func printNBTJSON() {
	/*remainingBytes, err := io.ReadAll(r)
	if err != nil {
		log.Fatal(err)
	}

	j, err := nbt2json.Nbt2Json(remainingBytes, "")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(j))*/
}
