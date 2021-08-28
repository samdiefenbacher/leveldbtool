package world

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/danhale-git/mine/mock"
)

var testWorld *World

const worldDirName = `97caYQjdAgA=`

func init() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("getting working directory: %s", err)
	}

	// Benchmarks need to be run from current dir
	if !strings.HasSuffix(wd, "world") {
		return
	}

	testWorld, err = New(filepath.Join(wd, worldDirName))
	if err != nil {
		log.Fatalf("unexpected error opening world: %s", err)
	}
}

var result Block

func BenchmarkGetBlock(b *testing.B) {
	if testWorld == nil {
		fmt.Println("test world is nil, are you in the world package directory?")
	}

	var r Block
	var err error

	for n := 0; n < b.N; n++ {
		r, err = testWorld.GetBlock(0, 0, 0, 0)
		if err != nil {
			b.Errorf("error returned getting block")
		}
	}

	result = r
}

func TestGetBlock(t *testing.T) {
	w := World{
		db:        mock.ValidLevelDB(),
		subChunks: make(map[struct{ x, y, z, d int }]*subChunkData),
	}

	expected := []Block{
		{Y: 0, id: "minecraft:crimson_planks", waterLogged: false, X: 0, Z: 0},
		{Y: 1, id: "minecraft:fence", waterLogged: true, X: 0, Z: 0},
		{Y: 2, id: "minecraft:air", waterLogged: false, X: 0, Z: 0},
	}

	for y := 0; y < 3; y++ {
		b, err := w.GetBlock(0, y, 0, 0)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		if b != expected[y] {
			t.Errorf("block did not match expected values: expected %+v: got %+v", expected[y], b)
		}
	}
}
