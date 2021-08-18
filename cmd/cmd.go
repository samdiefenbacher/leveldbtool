package cmd

import (
	"encoding/json"
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

//const worldFileName = `VsgSYaaGAAA=` // MINETEST  16 64 16
const worldFileName = `97caYQjdAgA=` // MINETESTFLAT 0 0 0

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

			//fmt.Println(mock.ByteSliceAsString(data))

			sc, err := terrain.NewSubChunk(data)
			if err != nil {
				log.Fatal(err)
			}

			prettyPrint(sc.Blocks.Palette)

			printCount := 48

			for i := range sc.Blocks.Indices {
				if i > printCount {
					break
				}

				idx := sc.Blocks.Indices[i]
				if idx >= len(sc.Blocks.Palette) {
					log.Printf("index %d out of range %d", idx, len(sc.Blocks.Palette))
					continue
				}
				block := sc.Blocks.Palette[idx].BlockID()

				var waterLogged string
				if len(sc.WaterLogged.Indices) > 0 {
					waterLogged = sc.WaterLogged.Palette[sc.WaterLogged.Indices[i]].BlockID()
				}

				x, y, z := blockPosition(i)

				fmt.Printf("(%d, %d, %d)", x, y, z)
				fmt.Printf(" - %s - %s [%d] (%d) \n", block, waterLogged, i, idx)
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

func prettyPrint(s interface{}) {
	b, err := json.MarshalIndent(s, "", "  ")
	fmt.Println(string(b), err)
}

func printAlignedSlice(a, b []int) {
	for i := range a {
		fmt.Printf("%d,%d\t", a[i], b[i])
	}
	fmt.Println()
}
