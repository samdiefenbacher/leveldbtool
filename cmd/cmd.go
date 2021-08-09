package cmd

import (
	"fmt"
	"log"
	"strconv"

	"github.com/danhale-git/mine/parse"

	"github.com/danhale-git/mine/terrain"
	"github.com/midnightfreddie/McpeTool/world"
	"github.com/spf13/cobra"
)

const worldPath = `C:\Users\danha\AppData\Local\Packages\Microsoft.MinecraftUWP_8wekyb3d8bbwe\LocalState\games\com.mojang\minecraftWorlds\XfILYVNgAQA=`

func Init() error {
	root := &cobra.Command{
		Use: "mine <x> <y> <z>",
		Run: func(cmd *cobra.Command, args []string) {
			w, err := world.OpenWorld(worldPath)
			if err != nil {
				log.Fatal(err)
			}
			defer w.Close()

			key, err := parse.SubChunkKey(
				int32(intArg(args[0])), // x
				int32(intArg(args[2])), // z
				0,                      // dimension
				intArg(args[1]),        // y
			)
			if err != nil {
				log.Fatal(err)
			}

			value, err := w.Get(key)
			if err != nil {
				log.Fatal(err)
			}
			sc, err := terrain.NewSubChunk(value)
			if err != nil {
				log.Fatal(err)
			}

			for i := range sc.BlockStorage {
				if i > 20 {
					break
				}
				index := sc.BlockStorage[i]
				x, y, z := blockPosition(i)

				fmt.Printf("(%d, %d, %d)", x, y, z)
				fmt.Printf(" - %s\n", sc.StatePalette[index].BlockID())
			}

			//util.ParseSubChunk(value)
		},
	}

	return root.Execute()
}

func intArg(a string) int {
	c, err := strconv.Atoi(a)
	if err != nil {
		log.Fatalf("invalid arg '%s': cannot convert to int: %s", a, err)
	}

	return c
}

// blockPosition gives the x/y/z offset from the a subchunk root based on the current index of the block storage record
func blockPosition(increment int) (x, y, z int) {
	x = (increment >> 8) & 0xF
	y = increment & 0xF
	z = (increment >> 4) & 0xF

	return
}
