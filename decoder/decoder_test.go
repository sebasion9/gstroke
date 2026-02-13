package decoder_test

import (
	"fmt"
	"gstroke/decoder"
	"slices"
	"testing"
)

// codes, lenght, maxbits
func TestCanonicalOk(t *testing.T) {
	tests := []struct {
		counts [16]uint8
		// expected
		maxBits int
		lenghts []int
		codes []int
	} {
		{
			counts: [16]uint8{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			maxBits: 1,
			lenghts: []int{1},
			codes: []int{0},
		},
		{
			counts: [16]uint8{2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			maxBits: 1,
			lenghts: []int{1, 1},
			codes: []int{0, 1},
		},
		{
			counts: [16]uint8{0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			maxBits: 2,
			lenghts: []int{2, 2},
			codes: []int{0, 1}, // 00, 01
		},
		{
			counts: [16]uint8{1, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			maxBits: 2,
			lenghts: []int{1, 2, 2},
			codes: []int{0, 2, 3}, // 0, 10, 11
		},
		{
			counts: [16]uint8{0, 1, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			maxBits: 3,
			lenghts: []int{2, 3, 3},
			codes: []int{0, 2, 3}, // 00, 010, 011
		},
		{
			counts: [16]uint8{1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			maxBits: 4,
			lenghts: []int{1, 2, 3, 4},
			codes: []int{0, 2, 6, 14}, // 0, 10, 110, 1110
		},
		{
			counts: [16]uint8{0, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			maxBits: 3,
			lenghts: []int{3, 3, 3, 3},
			codes: []int{0, 1, 2, 3}, // 000, 001, 010, 011
		},
		{
			counts: [16]uint8{0, 1, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			maxBits: 4,
			lenghts: []int{2, 4, 4},
			codes: []int{0, 4, 5}, // 00, 0100, 0101
		},
		{
			counts: [16]uint8{0, 2, 3, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			maxBits: 4,
			lenghts: []int{2, 2, 3, 3, 3, 4},
			codes: []int{0, 1, 4, 5, 6, 14}, // 00, 01, 100, 101, 110, 1110
		},
		{
			counts: [16]uint8{0, 1, 2, 3, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			maxBits: 5,
			lenghts: []int{2, 3, 3, 4, 4, 4, 5},
			codes: []int{0, 2, 3, 8, 9, 10, 22}, // 00, 010, 011, 1000, 1001, 1010, 10110
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("canonical: %d", i), func(t *testing.T) {
			ht := decoder.HuffTable{Counts: tt.counts}
			symbolsNum := 0
			for _, c := range tt.counts {
				symbolsNum += int(c)
			}
			ht.Symbols = make([]uint8, symbolsNum)

			ht.BuildCanonical()

			if !slices.Equal(ht.Codes, tt.codes) {
				t.Fatalf("codes arrays are not equal:\ngot: %v\nwant: %v", ht.Codes, tt.codes)
			}

			if !slices.Equal(ht.Lengths, tt.lenghts) {
				t.Fatalf("lengths arrays are not equal:\ngot: %v\nwant: %v", ht.Lengths, tt.lenghts)
			}

			if ht.MaxBits != tt.maxBits {
				t.Fatalf("maxBits are not equal:\ngot: %d\nwant: %d", ht.MaxBits, tt.maxBits)
			}
		})
	}

}
