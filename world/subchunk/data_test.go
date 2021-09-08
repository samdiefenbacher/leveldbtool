package subchunk

import "testing"

func TestBlockState(t *testing.T) {

}

func TestVoxelToIndex(t *testing.T) {
	i := 0
	for x := 0; x < 16; x++ {
		for z := 0; z < 16; z++ {
			for y := 0; y < 16; y++ {
				converted := voxelToIndex(x, y, z)
				if converted != i {
					t.Fatalf("expected coordinate %d, %d, %d to have index %d but got: %d",
						x, y, z, i, converted)
				}
				i++
			}
		}
	}
}

func TestIndexToVoxel(t *testing.T) {
	i := 0
	for x := 0; x < 16; x++ {
		for z := 0; z < 16; z++ {
			for y := 0; y < 16; y++ {
				cx, cy, cz := indexToVoxel(i)
				if cx != x || cy != y || cz != z {
					t.Fatalf("expected index %d to have coordinate %d %d %d but got: %d %d %d",
						i, x, y, z, cx, cy, cz)
				}
				i++
			}
		}
	}
}
