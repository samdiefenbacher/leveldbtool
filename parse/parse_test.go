package parse

import (
	"testing"

	"github.com/danhale-git/mine/mock"
)

func TestStorageCount(t *testing.T) {
	r := mock.SubChunkReader()
	l := r.Len()
	count, err := StorageCount(r)

	if err != nil {
		t.Errorf("unexpected error returned")
	}

	if count != mock.StorageCount {
		t.Errorf("unexpected storage count %d: expected %d", count, mock.StorageCount)
	}

	if r.Len()+2 != l {
		t.Errorf("this function should read 2 bytes but length changed from %d to %d", l, r.Len())
	}
}

func TestBlockStateIndices(t *testing.T) {
	r := mock.SubChunkReader()
	_, _ = r.Read(make([]byte, 2))

	indices, err := BlockStateIndices(r)
	if err != nil {
		t.Errorf("unexpected error returned: %s", err)
	}

	if len(indices) != subChunkBlockCount {
		t.Errorf("expected %d blocks state indices: got %d", subChunkBlockCount, len(indices))
	}
}
