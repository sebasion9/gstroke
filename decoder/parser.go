package decoder

import (
	"encoding/binary"
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

type Parser struct {
	source []byte
	pos int
	endPos int
	sos []byte
}

func newParser(source []byte) *Parser {
	return &Parser{
		source: source, 
	}
}


func (p* Parser) searchSeg(seg validSeg) int {
	srcLen := len(p.source)
	i := p.pos
	for ; i < srcLen - 1; i++ {
		marker := uint16(p.source[i]) << 8 | uint16(p.source[i+1])
		if validSeg(marker) == seg {
			return i
		}

	}

	return -1
}

// define quantisation table, len 69 bytes
func (p* Parser) parseDQT() ([]QuantTable, error) {
	var tables []QuantTable

	// advance marker bytes
	pos := p.pos + 2
	if pos + 2 > len(p.source) {
		return nil, errors.NewInvalidJPEGError("Invalid DQT segment")
	}

	segSize := int(binary.BigEndian.Uint16(p.source[pos:pos+2]))
	// advance seg size
	pos += 2
	// end of segment
	end := pos + (segSize - 2)


	for ; pos < end; {
		header := p.source[pos]
		precision := header >> 4
		id := header & 0x0F

		// advance dqt header
		pos++

		bytesPerVal := 1
		if precision == 1 {
			bytesPerVal = 2
		} else if precision != 0 {
			return nil, errors.NewInvalidJPEGError("Invalid DQT precision")
		}

		size := 64 * bytesPerVal

		if pos + size > end {
			return nil, errors.NewInvalidJPEGError("Truncated DQT table")
		}

		data := p.source[pos:pos+size]
		pos += size

		var values [64]uint16
		for i := range 64 {
			if precision == 0 {
				values[i] = uint16(data[i])
				continue
			}
			values[i] = binary.BigEndian.Uint16(data[i*2:i*2+2])
		}

		tables = append(tables, QuantTable{
			Precision: precision,
			ID: id,
			values: values,
		})

	}


	p.pos = end


	return tables, nil
}

// define huffman table
func (p* Parser) parseDHT() ([]HuffTable, error) {
	var tables []HuffTable
	// advance marker bytes
	pos := p.pos + 2
	if pos + 2 > len(p.source) {
		return nil, errors.NewInvalidJPEGError("Invalid DHT segment")
	}

	segSize := int(binary.BigEndian.Uint16(p.source[pos:pos+2]))
	pos += 2
	end := pos + segSize - 2

	for ; pos < end; {
		header := p.source[pos]
		pos++

		class := header >> 4
		id := header & 0x0F
		if class > 1 {
			return nil, errors.NewInvalidJPEGError("Invalid DHT class")
		}

		var counts [16]uint8
		copy(counts[:], p.source[pos:pos+16])
		pos+=16

		symbolsNum := 0
		for _, c := range counts {
			symbolsNum += int(c)
		}

		if pos + symbolsNum > end {
			return nil, errors.NewInvalidJPEGError("Truncated DHT symbols")
		}

		symbols := make([]uint8, symbolsNum)
		copy(symbols, p.source[pos:pos+symbolsNum])
		pos += symbolsNum


		tables = append(tables, HuffTable{
			Class: huffClass(class),
			ID: id,
			Counts: counts,
			Symbols: symbols,
		})

	}


	p.pos = end

	return tables, nil
}

// start of frame, entropy-coded baseline frame
func(p Parser) parseSOF() (StartOfFrame, error) {
	sof := StartOfFrame{}
	// advance marker bytes
	pos := p.pos + 2
	if pos + 2 > len(p.source) {
		return sof, errors.NewInvalidJPEGError("Invalid SOF segment")
	}

	segSize := int(binary.BigEndian.Uint16(p.source[pos:pos+2]))
	pos += 2
	end := pos + segSize - 2

	if pos + 4 > end {
		return sof, errors.NewInvalidJPEGError("Truncated SOF segment")
	}

	precision := p.source[pos]
	if precision != 8 {
		return sof, errors.NewInvalidJPEGError("SOF segment invalid, only supports precision = 8")
	}

	pos++
	y := binary.BigEndian.Uint16(p.source[pos:pos+2])
	x := binary.BigEndian.Uint16(p.source[pos+2:pos+4])
	pos += 4

	nf := uint8(p.source[pos])
	pos++

	var components []Component
	for range nf {
		if pos + 2 > end {
			return sof, errors.NewInvalidJPEGError("Truncated SOF segment")
		}
		cid := p.source[pos]
		hv := p.source[pos+1]
		tq := p.source[pos+2]
		if tq > 3 {
			return sof, errors.NewInvalidJPEGError("Invalid Tqi value for component")
		}

		h := hv >> 4
		v := hv & 0x0f

		components = append(components, Component{
			H: h,
			V: v,
			CID: cid,
			Tq: tq,
		})

		pos += 3
	}

	sof.Components = components
	sof.Nf = nf
	sof.Precision = precision
	sof.X = x
	sof.Y = y


	p.pos = pos
	return sof, nil
}

//TODO:
// start of scan
func(p Parser) parseSOS() (StartOfScan, error) {
	sos := StartOfScan{}
	// advance marker bytes
	pos := p.pos + 2
	if pos + 2 > len(p.source) {
		return sos, errors.NewInvalidJPEGError("Invalid SOF segment")
	}

	segSize := int(binary.BigEndian.Uint16(p.source[pos:pos+2]))
	pos += 2
	end := pos + segSize - 2

	ns := uint8(p.source[pos])
	pos++
	var scanComponents []ScanComponent

	for range ns {
		if pos + 1 > end {
			return sos, errors.NewInvalidJPEGError("Truncated SOS segment")
		}

		cs := p.source[pos]
		tx := p.source[pos+1]

		td := tx >> 4
		ta := tx & 0x0f

		scanComponents = append(scanComponents, ScanComponent{
			Cs: cs,
			Td: td,
			Ta: ta,
		})


		pos += 2
	}

	if pos + 3 > end {
		return sos, errors.NewInvalidJPEGError("Truncated SOS segment")
	}

	ss := uint8(p.source[pos])
	se := uint8(p.source[pos+1])
	succApprox := p.source[pos+2]
	sah := succApprox >> 4
	sal := succApprox & 0x0f

	sos.ScanComponents = scanComponents
	sos.Ns = ns
	sos.SuccApproxH = sah
	sos.SuccApproxL = sal
	sos.StartSpectralPredictor = ss
	sos.EndSpectralPredictor = se

	return sos, nil
}

//TODO: verify scan works correctly
func(p *Parser) parseScan() ([]byte, error) {
	if p.pos >= p.endPos {
		return nil, errors.NewInvalidJPEGError("Invalid scan size")
	}
	var table []byte
	ff := false
	for i := p.pos; i < p.endPos; i++ {
		b := p.source[i]
		if b == 0xFF {
			ff = true
			continue
		}
		if ff {
			if b == 0x00 {
				table = append(table, 0xFF)
				ff = false
			}
			continue
		}
		table = append(table, b)
	}

	p.pos = p.endPos
	return table, nil
}
