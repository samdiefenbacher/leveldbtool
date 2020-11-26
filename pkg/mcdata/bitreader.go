package mcdata

import (
	"bytes"
)

const bitsPerByte = 8

// BitReader ready a bytes.Reader bit by bit.
type BitReader struct {
	reader *bytes.Reader // Next bytes
	byte   byte          // Current byte
	offset byte          // Offset from the previous byte
}

func NewBitReader(reader *bytes.Reader) BitReader {
	return BitReader{reader: reader}
}

// Offset returns the current offset from the previous byte, in bits.
func (r *BitReader) Offset() int {
	return int(r.offset)
}

// ReadBits reads the next bits from stored bytes and returns them.
func (r *BitReader) ReadBits(count int) ([]bool, error) {
	b := make([]bool, count)

	for i := 0; i < count; i++ {
		bit, err := r.ReadBit()
		if err != nil {
			return b, err
		}

		b[i] = bit
	}

	return b, nil
}

// ReadBit reads the next bit from stored bytes and returns it.
func (r *BitReader) ReadBit() (bool, error) {
	if r.offset == bitsPerByte {
		r.offset = 0
	}

	if r.offset == 0 {
		var err error
		if r.byte, err = r.reader.ReadByte(); err != nil {
			return false, err
		}
	}

	bit := (r.byte & (0x80 >> r.offset)) != 0

	r.offset++

	return bit, nil
}
