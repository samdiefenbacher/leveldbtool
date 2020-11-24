package main

import (
	"fmt"
	"log"

	"github.com/danhale-git/mine/pkg/mcdata"

	"github.com/danhale-git/mine/pkg/convert"

	"github.com/midnightfreddie/McpeTool/world"
	//"github.com/midnightfreddie/nbt2json"
)

//https://github.com/midnightfreddie/McpeTool/blob/master/examples/PowerShell/CsCoords.ps1
//https://minecraft.gamepedia.com/Bedrock_Edition_level_format

// https://github.com/midnightfreddie/McpeTool/tree/master/docs#how-to-convert-world-coordinates-to-leveldb-keys

func init() {
	log.SetFlags(log.Lshortfile)
}

func main() {
	/*path := "C:\\Users\\danha\\AppData\\Local\\Packages\\Microsoft.MinecraftUWP_8wekyb3d8bbwe\\LocalState\\games" +
	"\\com.mojang\\minecraftWorlds\\qOV5X3kvAAA="*/
	path := "C:\\Users\\danha\\AppData\\Local\\Packages\\Microsoft.MinecraftUWP_8wekyb3d8bbwe\\LocalState\\games" +
		"\\com.mojang\\minecraftWorlds\\4xq8X8xLAAA="
	w, err := world.OpenWorld(path)

	if err != nil {
		log.Println("error opening world", err)
	}

	/*b, err := w.GetKeys()

	if err != nil {
		log.Fatal(err)
	}

	for _, key := range b[:50] {
		dst := make([]byte, 0)
		fmt.Println(hex.Decode(dst, key))
	}
	*/

	x, y, z := 1, 40, 1

	key := convert.CoordsToSubChunkKey(x, y, z, 0)

	log.Printf("CoordsToSubChunkKey: %x", key)

	b, err := w.Get(key)

	if err != nil {
		log.Println(err)
	}

	mcdata.NewSubChunk(b)

	err = w.Close()

	if err != nil {
		fmt.Println(err)
	}
}
