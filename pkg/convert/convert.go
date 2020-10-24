package convert

import (
	"encoding/binary"
	"log"
)

func CoordsToDBKey(x, y, z, dimension int) string {
	xHex := convertInt(uint32(x / 16))
	zHex := convertInt(uint32(z / 16))

	log.Printf("%d, %d - %x, %x", x, y, xHex, zHex)

	return ""
}

func convertInt(i uint32) []byte {
	b := make([]byte, 32)
	binary.LittleEndian.PutUint32(b, i)
	return b[:4]
}
