package lazyledger

import (
    "hash"
    "math/bits"
)

// sum returns the hash of the input data using the specified algorithm.
func sum(h hash.Hash, data ...[]byte) []byte {
	h.Reset()
	for _, d := range data {
		// the Hash interface specifies that Write never returns an error
		_, _ = h.Write(d)
	}
	return h.Sum(nil)
}

// leafSum returns the hash created from data inserted to form a leaf. Leaf
// sums are calculated using:
//		Hash(0x00 || data)
func leafSum(h hash.Hash, data []byte) []byte {
	return sum(h, []byte{0x00}, data)
}

// nextSubtreeSize returns the size of the subtree adjacent to start that does
// not overlap end.
func nextSubtreeSize(start, end uint64) int {
	ideal := bits.TrailingZeros64(start)
	max := bits.Len64(end-start) - 1
	if ideal > max {
		return 1 << uint(max)
	}
	return 1 << uint(ideal)
}
