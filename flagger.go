package lazyledger

// Flagger is an interface for computing the flags of bytes of data.
type Flagger interface {
    // LeafFlag returns the flag of a raw unhashed leaf.
    LeafFlag([]byte) []byte

    // NodeFlag returns the flag of an intermediate node.
    NodeFlag([]byte) []byte

    // Union returns the union of two flags.
    Union([]byte, []byte) []byte

    // FlagSize returns the number of bytes that LeafFlag or Union will return.
    FlagSize() int
}
