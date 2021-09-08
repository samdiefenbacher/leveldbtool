package subchunk

import (
	"testing"

	"github.com/danhale-git/mine/mock"
)

func TestNew(t *testing.T) {
	_, err := Decode(mock.SubChunkValue)
	if err != nil {
		t.Errorf("unexpected error returned: %s", err)
	}
}

func TestStateIndices(t *testing.T) {
	r := mock.SubChunkReader()
	_, _ = r.Read(make([]byte, 2))

	indices, err := stateIndices(r)
	if err != nil {
		t.Errorf("unexpected error returned: %s", err)
	}

	if len(indices) != BlockCount {
		t.Errorf("expected %d blocks state indices: got %d", BlockCount, len(indices))
	}
}
