package cmd

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/danhale-git/mine/util"
	"github.com/midnightfreddie/McpeTool/world"
	"github.com/spf13/cobra"
	"io"
	"log"
	"strconv"
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

			key, err := util.SubChunkKey(
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

			r := bytes.NewReader(value)

			var version int8
			if err := read(r, &version); err != nil {
				log.Fatal(err)
			}

			var storageCount int8 = 1
			if version != 1 {
				if err := read(r, &storageCount); err != nil {
					log.Fatal(err)
				}
			}

			fmt.Println("version:", version)
			fmt.Println("storageCount:", storageCount)
		},
	}

	return root.Execute()
}

func read(r io.Reader, data interface{}) error {
	return binary.Read(r, binary.ByteOrder(binary.LittleEndian), data)
}

func intArg(a string) int {
	c, err := strconv.Atoi(a)
	if err != nil {
		log.Fatalf("invalid arg '%s': cannot convert to int: %s", a, err)
	}

	return c
}
