package main

import (
	"fmt"
	"log"

	"github.com/danhale-git/mine/pkg/mcdata"

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

	// Determine key from XYZ coordinates
	x, _, z := 1, 40, 1

	wld := mcdata.World{World: &w}

	c, err := wld.Chunk(x, z)

	if err != nil {
		log.Fatal(err)
	}

	for _, sc := range c.SubChunks {
		for _, bs := range sc.BlockStorage {
			PrintBlockStorage(bs)
			fmt.Println("===============")
		}
	}

	err = w.Close()

	if err != nil {
		fmt.Println(err)
	}
}

func PrintBlockStorage(blocks mcdata.BlockStorage) {
	uniqueCounts := make(map[string]int)

	for _, idx := range blocks.BlockStateIndices {
		description := ""

		name, err := blocks.BlockName(idx)
		description += name

		if err != nil {
			log.Fatal(err)
		}

		states, err := blocks.BlockStateTags(idx)
		if err != nil {
			log.Fatal(err)
		}

		for _, state := range states {
			description += fmt.Sprintf(" - %v %v", state.Name, state.Value)
		}

		if _, ok := uniqueCounts[description]; !ok {
			uniqueCounts[description] = 1
		} else {
			uniqueCounts[description]++
		}
	}

	total := 0
	for k, v := range uniqueCounts {
		fmt.Println(k, v)
		total += v
	}

	fmt.Println("total blocks -", total)
}
