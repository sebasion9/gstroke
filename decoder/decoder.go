package decoder

import (
	"fmt"
	"gstroke/errors"
)

type QuantTable struct {
	ID uint8 //  Tq
	Precision uint8 // Pq -> 8 or 16bit
	values [64]uint16
}

type huffClass uint8
const (
	DC huffClass = 0
	AC huffClass = 1
)

type HuffTable struct {
	Class huffClass
	ID uint8
	Counts [16]uint8
	Symbols []uint8
}

type Decoder struct {
	*Parser
	dqt []QuantTable
	dht []HuffTable
}

func NewDecoder(source []byte) *Decoder {
	return &Decoder{
		Parser: newParser(source),
	}
}

func (d* Decoder) Decode() error {
	fmt.Println("[INFO] start parsing (delete later)")

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

		tables, err := d.parseDQT()
		if err != nil {
			return err
		}

		d.dqt = append(d.dqt, tables...)

	}

	for {
		pos = d.searchSeg(DHT)
		if pos == -1 && len(d.dht) == 0 {
			return errors.NewInvalidJPEGError("No DHT markers found")
		} else if (pos == -1) {
			break
		}
		d.pos = pos


		tables, err := d.parseDHT()
		if err != nil {
			return err
		}

		d.dht = append(d.dht, tables...)

	}
	// restart positon to 0
	d.pos = 0

	d.pos = d.searchSeg(SOF)
	if d.pos == -1 {
		return errors.NewInvalidJPEGError("No SOF marker")
	}

	//TODO:
	if err := d.parseSOF(); err != nil { return err }

	d.pos = d.searchSeg(SOS)
	if d.pos == -1 {
		return errors.NewInvalidJPEGError("No SOS marker")
	}


	//TODO:
	if err := d.parseSOS(); err != nil { return err }
	if err := d.parseScan(); err != nil { return err }

	return nil
}
