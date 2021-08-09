package main

// replicate `.\getdata_rawcoords.ps1 1 1 1`

import (
	"github.com/danhale-git/mine/cmd"
	"log"
)

func main() {
	if err := cmd.Init(); err != nil {
		log.Fatal(err)
	}
}
