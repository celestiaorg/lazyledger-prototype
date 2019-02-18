package lazyledger

type Flagger interface {
    // GetLeafFlag returns the flag of a leaf node.
    GetLeafFlag([]byte) []byte

    // Union returns the union of two flags.
    Union([]byte, []byte)

    // Size returns the fixed size of all flags.
    Size() int
}
