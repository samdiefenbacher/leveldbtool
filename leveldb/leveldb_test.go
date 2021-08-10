package leveldb

import (
	"encoding/hex"
	"strings"
	"testing"
)

func TestSubChunkKey(t *testing.T) {
	testSubChunkKey(0, 0, 0, "00000000000000002F00", t)
	testSubChunkKey(16, 16, 16, "01000000010000002F01", t)
	testSubChunkKey(-1, 32, -1, "FFFFFFFFFFFFFFFF2F02", t)
}

func testSubChunkKey(x, y, z int, want string, t *testing.T) {
	b, err := SubChunkKey(x, y, z, 0)
	if err != nil {
		t.Errorf("unexpected error returned: %s", err)
	}

	got := strings.ToUpper(hex.EncodeToString(b))

	if want != got {
		t.Errorf("unexpected key '%s': expected '%s'", got, want)
	}
}
