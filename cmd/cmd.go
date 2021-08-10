package cmd

import (
	"fmt"
	"log"
	"path/filepath"
	"strconv"

	"github.com/midnightfreddie/McpeTool/world"

	"github.com/danhale-git/mine/leveldb"

	"github.com/danhale-git/mine/terrain"
	"github.com/spf13/cobra"
)

const worldDirPath = `C:\Users\danha\AppData\Local\Packages\Microsoft.MinecraftUWP_8wekyb3d8bbwe\LocalState\games\com.mojang\minecraftWorlds\`
const worldFileName = `VsgSYaaGAAA=`

func Init() error {
	root := &cobra.Command{
		Use: "mine <x> <y> <z>",
		Run: func(cmd *cobra.Command, args []string) {
			w, err := world.OpenWorld(filepath.Join(worldDirPath, worldFileName))
			if err != nil {
				log.Fatal(err)
			}

			key, err := leveldb.SubChunkKey(
				intArg(args[0]), // x
				intArg(args[1]), // y
				intArg(args[2]), // z
				0,
			)

			data, err := w.Get(key)
			if err != nil {
				log.Fatal(err)
			}

			sc, err := terrain.NewSubChunk(data)
			if err != nil {
				log.Fatal(err)
			}

			for i := range sc.BlockStorage {
				if i > 32 {
					break
				}
				index := sc.BlockStorage[i]
				x, y, z := blockPosition(i)

				fmt.Printf("(%d, %d, %d)", x, y, z)
				fmt.Printf(" - %s\n", sc.StatePalette[index].BlockID())
			}
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
