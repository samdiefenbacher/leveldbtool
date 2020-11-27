package mcdata

type Block struct {
	x, y, z int
	Name    string
	States  map[string]interface{}
	Version int
}
