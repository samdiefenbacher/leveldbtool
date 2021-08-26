package world

import (
	"testing"

	"github.com/danhale-git/mine/mock"
)

func TestGetBlock(t *testing.T) {
	w := World{mock.ValidLevelDB()}

	expected := []Block{
		{Y: 0, id: "minecraft:crimson_planks", waterLogged: false, X: 0, Z: 0},
		{Y: 1, id: "minecraft:fence", waterLogged: true, X: 0, Z: 0},
		{Y: 2, id: "minecraft:air", waterLogged: false, X: 0, Z: 0},
	}

	for y := 0; y < 2; y++ {
		b, err := w.GetBlock(0, y, 0, 0)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		if b != expected[y] {
			t.Errorf("block did not match expected values: expected %+v: got %+v", expected[y], b)
		}
	}
}
