package main

import (
	"log"

	"github.com/danhale-git/mine/pkg/convert"

	"github.com/midnightfreddie/McpeTool/world"
)

const blocksPerChunk = 16.0

//https://github.com/midnightfreddie/McpeTool/blob/master/examples/PowerShell/CsCoords.ps1
//https://minecraft.gamepedia.com/Bedrock_Edition_level_format

// https://github.com/midnightfreddie/McpeTool/tree/master/docs#how-to-convert-world-coordinates-to-leveldb-keys

func init() {
	log.SetFlags(log.Lshortfile)
}

func main() {
	path := "C:\\Users\\danha\\AppData\\Local\\Packages\\Microsoft.MinecraftUWP_8wekyb3d8bbwe\\LocalState\\games\\com.mojang\\minecraftWorlds\\qOV5X3kvAAA="
	_, err := world.OpenWorld(path)

	if err != nil {
		log.Println("error opening world", err)
	}

	convert.CoordsToDBKey(413, 105, 54, 0)

	if err != nil {
		log.Println(err)
	}

	/*b, err := w.Get([]byte(GetKeyByCoords(17, 17, 50, 0)))
	log.Println(b, err)
	err = w.Close()
	if err != nil {
		fmt.Println(err)
	}*/
}

/*func GetKeyByCoords(x, z, y, dimension int) string {
	var Tag byte = 0x2f

	var d string
	if dimension != 0 {
		d = hex(dimension)
	}

	var t string

	if Tag == byte(0x2f) {
		SubChunkY := byte(y / blocksPerChunk)
		t = hex(SubChunkY)
	}

	ChunkX := int32(math.Floor(float64(x) / blocksPerChunk))
	ChunkZ := int32(math.Floor(float64(z) / blocksPerChunk))

	log.Println("val:", ChunkX, "hex:", hex(ChunkX))

	return hex(ChunkX) + hex(ChunkZ) + d + hex(Tag) + t
}

func hex(i interface{}) string {
	return fmt.Sprintf("%x", i)
}*/
