package mock

type LevelDB struct {
	data []byte
}

func (w *LevelDB) Get(_ []byte) ([]byte, error) {
	return w.data, nil
}

func ValidLevelDB() *LevelDB {
	return &LevelDB{SubChunkValue}
}
