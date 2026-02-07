package decoder

import "fmt"

func (h *HuffTable) BuildCanonical() {
	h.Codes = make([]int, len(h.Symbols))
	h.Lengths = make([]int, len(h.Symbols))

	code := 0
	symbolIdx := 0
	maxBits := 0


	for length := 1; length <= 16; length++ {
		count := h.Counts[length - 1]
		if count > 0 {
			maxBits = int(length)
		}

		for i := 0; i < int(count); i++ {
			h.Codes[symbolIdx] = code
			h.Lengths[symbolIdx] = length
			code++
			symbolIdx++
		}
		code <<= 1
	}

	h.MaxBits = maxBits
	fmt.Printf("%+v\n", h)
}
