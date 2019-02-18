package lazyledger

import (
    "hash"
)

type flagDigest struct {
    flagger Flagger
    baseHasher hash.Hash
    leaf bool
    currentData []byte
}

func NewFlagHasher(flagger Flagger, baseHasher hash.Hash) hash.Hash {
    return &flagDigest{
        flagger: flagger,
        baseHasher: baseHasher,
    }
}

func (d *flagDigest) Write(p []byte) (int, error) {
    if d.currentData == nil {
        if p[0] == byte(0) {
            d.leaf = true
        } else {
            d.leaf = false
        }
    }

    d.currentData = append(d.currentData, p...)
    return d.baseHasher.Write(p)
}

func (d *flagDigest) Sum(in []byte) []byte {
    // TODO: add flags here
    return d.baseHasher.Sum(in)
}

func (d *flagDigest) Size() int {
    return d.flagger.Size() + d.baseHasher.Size()
}

func (d *flagDigest) BlockSize() int {
    return d.baseHasher.BlockSize()
}

func (d *flagDigest) Reset() {
    d.currentData = nil
    d.baseHasher.Reset()
}
