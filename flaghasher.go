package lazyledger

import (
    "hash"
)

type flagDigest struct {
    flagger Flagger
    baseHasher hash.Hash
    data []byte
}

// NewFlagHasher returns a new hash.Hash computing checksums using the bashHasher with flags from flagger.
func NewFlagHasher(flagger Flagger, baseHasher hash.Hash) hash.Hash {
    return &flagDigest{
        flagger: flagger,
        baseHasher: baseHasher,
    }
}

func (d *flagDigest) Write(p []byte) (int, error) {
    d.data = append(d.data, p...)
    return d.baseHasher.Write(p)
}

func (d *flagDigest) Sum(in []byte) []byte {
    in = append(in, d.parentFlag()...)
    return d.baseHasher.Sum(in)
}

func (d *flagDigest) Size() int {
    return d.flagger.FlagSize() + d.baseHasher.Size()
}

func (d *flagDigest) BlockSize() int {
    return d.baseHasher.BlockSize()
}

func (d *flagDigest) Reset() {
    d.data = nil
    d.baseHasher.Reset()
}

func (d *flagDigest) leftFlag() []byte {
    return d.flagger.NodeFlag(d.data[1:d.Size()+1])
}

func (d *flagDigest) rightFlag() []byte {
    return d.flagger.NodeFlag(d.data[1+d.Size():])
}

func (d *flagDigest) parentFlag() []byte {
    if d.isLeaf() {
        return d.flagger.LeafFlag(d.mainData())
    }
    return d.flagger.Union(d.leftFlag(), d.rightFlag())
}

func (d *flagDigest) mainData() []byte {
    return d.data[1:]
}

func (d *flagDigest) isLeaf() bool {
    if d.data[0] == byte(0) {
        return true
    }
    return false
}
