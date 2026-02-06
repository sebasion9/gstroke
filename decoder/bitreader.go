package decoder

type BitReader struct {
	bitPos int
	bytePos int
	scan []byte
}

func newBitReader() *BitReader {
	return &BitReader{bitPos: 0, bytePos: 0, scan: []byte{}}
}

func (b *BitReader) ReadBit() int {
	bit := int((b.scan[b.bytePos] >> (7 - b.bitPos)) & 1)

	b.bitPos++
	if b.bitPos > 7 {
		b.bitPos = 0
		b.bytePos++
	}
	return bit
}

func(b *BitReader) ReadBits(n int) int {
	bits := 0
	for range n {
		bits = (bits << 1) | b.ReadBit()
	}
	return bits
}
