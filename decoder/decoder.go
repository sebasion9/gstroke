package decoder

import (
	"encoding/binary"
	"fmt"
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
	sof []byte
	sos []byte
	scan []byte
}

func NewDecoder(source []byte) *Decoder {
	return &Decoder{
		source: source, 
		pos: 0,
	}
}


func (d* Decoder) Decode() error {
	fmt.Println("[INFO] start decoding (delete later)")

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

		if err := d.readDQT(); err != nil { return err }

	}

	for {
		pos = d.searchSeg(DHT)
		if pos == -1 && len(d.dht) == 0 {
			return errors.NewInvalidJPEGError("No DHT markers found")
		} else if (pos == -1) {
			break
		}
		d.pos = pos

		if err := d.readDHT(); err != nil { return err }

	}

	// restart positon to 0
	d.pos = 0

	d.pos = d.searchSeg(SOF)
	if d.pos == -1 {
		return errors.NewInvalidJPEGError("No SOF marker")
	}

	if err := d.readSOF(); err != nil { return err }

	d.pos = d.searchSeg(SOS)
	if d.pos == -1 {
		return errors.NewInvalidJPEGError("No SOS marker")
	}


	if err := d.readSOS(); err != nil { return err }
	if err := d.readScan(); err != nil { return err }

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
	// advance marker bytes
	pos := d.pos + 2
	if pos + 2 > len(d.source) {
		return errors.NewInvalidJPEGError("Invalid DQT segment")
	}

	segSize := int(binary.BigEndian.Uint16(d.source[pos:pos+2]))


	table := make([]byte, segSize-2)
	copy(table, d.source[pos+2:pos+segSize])

	d.dqt = append(d.dqt, table)
	d.pos = pos+segSize

	return nil
}

// define huffman table
func (d* Decoder) readDHT() error {
	// advance marker bytes
	pos := d.pos + 2
	if pos + 2 > len(d.source) {
		return errors.NewInvalidJPEGError("Invalid DHT segment")
	}

	segSize := int(binary.BigEndian.Uint16(d.source[pos:pos+2]))

	table := make([]byte, segSize-2)
	copy(table, d.source[pos+2:pos+segSize])

	d.dht = append(d.dht, table)
	d.pos = pos+segSize

	return nil
}

// start of frame, entropy-coded baseline frame
func(d* Decoder) readSOF() error {
	// advance marker bytes
	pos := d.pos + 2
	if pos + 2 > len(d.source) {
		return errors.NewInvalidJPEGError("Invalid SOF segment")
	}

	segSize := int(binary.BigEndian.Uint16(d.source[pos:pos+2]))

	table := make([]byte, segSize - 2)
	copy(table, d.source[pos+2:pos+segSize])
	d.sof = table

	d.pos = pos+segSize
	return nil
}

// start of scan
func(d* Decoder) readSOS() error {
	// advance marker bytes
	pos := d.pos + 2
	if pos + 2 > len(d.source) {
		return errors.NewInvalidJPEGError("Invalid SOF segment")
	}

	segSize := int(binary.BigEndian.Uint16(d.source[pos:pos+2]))

	table := make([]byte, segSize - 2)
	copy(table, d.source[pos+2:pos+segSize])
	d.sos = table

	d.pos = pos+segSize
	return nil
}

func(d *Decoder) readScan() error {
	table := make([]byte, d.endPos - d.pos)
	copy(table, d.source[d.pos:d.endPos])
	d.scan = table
	d.pos = d.endPos
	return nil
}
