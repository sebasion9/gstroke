package decoder

import (
	"fmt"
	"gstroke/errors"
)

type huffClass uint8
const (
	DC huffClass = 0
	AC huffClass = 1
)

type QuantTable struct {
	ID uint8 //  Tq
	Precision uint8 // Pq -> 8 or 16bit
	values [64]uint16
}

type HuffTable struct {
	Class huffClass
	ID uint8
	Counts [16]uint8
	Symbols []uint8
}


type StartOfFrame struct {
	Precision uint8
	Y uint16
	X uint16
	Nf uint8
	Components []Component
}

type Component struct {
	CID uint8
	H uint8
	V uint8
	Tq uint8
}

type StartOfScan struct {
	Ns uint8
	ScanComponents []ScanComponent
	StartSpectralPredictor uint8
	EndSpectralPredictor uint8
	SuccApproxH uint8
	SuccApproxL uint8
}

type ScanComponent struct {
	Cs uint8
	Td uint8
	Ta uint8
}

type Decoder struct {
	*Parser
	*BitReader
	dqt []QuantTable
	dht []HuffTable
	sof StartOfFrame
	sos StartOfScan
	scan []byte
}

func NewDecoder(source []byte) *Decoder {
	return &Decoder{
		Parser: newParser(source),
		BitReader: newBitReader(),
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

	sof, err := d.parseSOF()
	if err != nil {
		return err
	}
	d.sof = sof

	d.pos = d.searchSeg(SOS)
	if d.pos == -1 {
		return errors.NewInvalidJPEGError("No SOS marker")
	}


	sos, err := d.parseSOS()
	if err != nil {
		return err
	}
	d.sos = sos

	scan, err := d.parseScan()
	if err != nil {
		return err
	}
	d.scan = scan
	d.BitReader.scan = scan

	d.decodeHuffman()

	return nil
}

func (d *Decoder) decodeHuffman() {

}

