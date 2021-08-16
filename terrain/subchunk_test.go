package terrain

import (
	"testing"

	"github.com/danhale-git/mine/mock"
)

func TestNewSubChunk(t *testing.T) {
	_, err := NewSubChunk(mock.SubChunkValue)
	if err != nil {
		t.Errorf("unexpected error returned: %s", err)
	}
}
