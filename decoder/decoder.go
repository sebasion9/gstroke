package decoder

import (
	"gstroke/errors"
)

type validSeg uint16
var (
	SOI validSeg = 0xFFD8
	DQT validSeg = 0xFFDB
	DHT validSeg = 0xFFC4
	SOF validSeg = 0xFFC0
	SOS validSeg = 0xFFDA
	EOI validSeg = 0xFFD9
)

type Decoder struct {
	source []byte
	pos int
	endPos int
	dqt [][]byte
	dht [][]byte
}

func NewDecoder(source []byte) *Decoder {
	return &Decoder{
		source: source, 
		pos: 0,
	}
}


func (d* Decoder) Decode() error {
	pos := d.searchSeg(SOI)
	if pos == -1 {
		return errors.NewInvalidJPEGError("No start marker")
	}

	d.endPos = d.searchSeg(EOI)
	if d.endPos == -1 {
		return errors.NewInvalidJPEGError("No end marker")
	}

	for {
		pos = d.searchSeg(DQT)
		if pos == -1 && len(d.dqt) == 0 {
			return errors.NewInvalidJPEGError("No DQT markers found")
		} else if (pos == -1) {
			break
		}
		d.pos = pos

		d.readDQT()
		// push to dht array
	}

	for {
		pos = d.searchSeg(DHT)
		if pos == -1 && len(d.dht) == 0 {
			return errors.NewInvalidJPEGError("No DHT markers found")
		} else if (pos == -1) {
			break
		}
		d.pos = pos

		d.readDHT()
		// push to dht array
	}

	d.pos = d.searchSeg(SOF)
	if d.pos == -1 {
		return errors.NewInvalidJPEGError("No SOF marker")
	}

	d.pos = d.searchSeg(SOS)
	if d.pos == -1 {
		return errors.NewInvalidJPEGError("No SOS marker")
	}


	return nil
}

func (d* Decoder) searchSeg(seg validSeg) int {
	srcLen := len(d.source)
	i := d.pos
	for ; i < srcLen - 1; i++ {
		marker := uint16(d.source[i]) << 8 | uint16(d.source[i+1])
		if validSeg(marker) == seg {
			return i
		}

	}

	return -1
}

// define quantisation table, len 69 bytes
func (d* Decoder) readDQT() error {

	// update d.pos
	return nil
}

// define huffman table
func (d* Decoder) readDHT() error {

	// update d.pos
	return nil
}
