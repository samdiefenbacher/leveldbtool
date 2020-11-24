package mcdata

import "bytes"

type BitReader struct {
	reader *bytes.Reader
	byte   byte
	offset byte
}

func NewBitReader(reader *bytes.Reader) BitReader {
	return BitReader{reader: reader}
}

func (r *BitReader) Offset() int {
	return int(r.offset)
}

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

func (r *BitReader) ReadBit() (bool, error) {
	if r.offset == 8 {
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
