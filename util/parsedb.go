package util

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"

	"github.com/midnightfreddie/nbt2json"
)

func read(r io.Reader, data interface{}) error {
	return binary.Read(r, binary.ByteOrder(binary.LittleEndian), data)
}

type NBTTag struct {
	Type  byte        `json:"tagType"`
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}

func (n *NBTTag) BlockID() string {
	//	fmt.Printf("%+v\n", n)
	if vs, ok := n.Value.([]interface{}); ok {
		for _, t := range vs {
			if tMap, ok := t.(map[string]interface{}); ok {
				if tMap["name"] == "name" {
					return tMap["value"].(string)
				}
			}
		}
	} else {
		fmt.Println("failed to convert to NBTTag")
	}

	return ""
}

func ParseSubChunk(data []byte) {
	//printNBTJSON(data)

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

	var s int8
	for s = 0; s < storageCount; s++ {
		var bitsPerBlockAndVersion byte
		if err := read(r, &bitsPerBlockAndVersion); err != nil {
			log.Fatal(err)
		}

		bitsPerBlock := int(bitsPerBlockAndVersion >> 1)
		storageVersion := int(bitsPerBlockAndVersion & 1)
		fmt.Println("bitsPerBlock:", float64(bitsPerBlock)) //DEBUG
		fmt.Println("storageVersion:", storageVersion)      //DEBUG

		//indices := make([]int, 0)

		indices, x, y, z := getBlockDataIndices(r, bitsPerBlock)

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

		remainingBytes, err := io.ReadAll(r)
		if err != nil {
			log.Fatal(err)
		}

		j, err := nbt2json.Nbt2Json(remainingBytes, "")
		if err != nil {
			log.Fatal(err)
		}

		nbt := struct {
			Nbt []NBTTag
		}{}
		//var nbtTags []NBTTag
		if err := json.Unmarshal(j, &nbt); err != nil {
			log.Fatal(err)
		}

		for i := range indices {
			if i > 5 {
				break
			}
			index := indices[i]
			x, y, z := x[i], y[i], z[i]

			fmt.Printf("(%d, %d, %d)", x, y, z)
			fmt.Printf(" - %s\n", nbt.Nbt[index].BlockID())
		}
	}
}

func getBlockDataIndices(r io.Reader, bitsPerBlock int) (indices, x, y, z []int) {
	blocksPerWord := int(math.Floor(32.0 / float64(bitsPerBlock)))

	fmt.Println("blocksPerWord:", blocksPerWord) //DEBUG

	wordCount := int(math.Ceil(4096 / float64(blocksPerWord)))
	fmt.Println("wordCount:", wordCount) //DEBUG

	indices = make([]int, 4096)
	x = make([]int, 4096)
	y = make([]int, 4096)
	z = make([]int, 4096)
	idx := 0

	for w := 0; w < wordCount; w++ {
		word := make([]byte, 4)
		if err := read(r, word); err != nil {
			log.Fatal(err)
		}

		// TODO: Also get the block position

		// Might need to use a bit reader here if numbers other than 4 or 8 come up
		switch bitsPerBlock {
		case 4:
			for _, b := range word {
				first := b >> 4
				x[idx], y[idx], z[idx] = blockPosition(idx)
				indices[idx] = int(first)
				idx++

				second := (b << 4) >> 4
				x[idx], y[idx], z[idx] = blockPosition(idx)
				indices[idx] = int(second)
				idx++
			}
		default:
			log.Panicf("unhandled bits per block '%d'", bitsPerBlock)
		}
	}

	return indices, x, y, z
}

func blockPosition(increment int) (x, y, z int) {
	x = (increment >> 8) & 0xF
	y = increment & 0xF
	z = (increment >> 4) & 0xF

	return
}

func printNBTJSON(b []byte) {
	j, err := nbt2json.Nbt2Json(b, "")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(j))
}
