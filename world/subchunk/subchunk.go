package subchunk

import (
	"math"
)

const BlockCount = 4096
const Size = 16
const WaterID = "minecraft:water"

// Origin returns the origin of the sub chunk containing the given coordinates. This is the corner block with
// the lowest x, y and z values.
func Origin(x, y, z int) (xo, yo, zo int) {
	xo = int(math.Floor(float64(x) / 16))
	yo = int(math.Floor(float64(y) / 16))
	zo = int(math.Floor(float64(z) / 16))

	return
}
