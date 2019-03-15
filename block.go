package lazyledger

// Block represents a block in the chain.
type Block interface {
    // AddMessage adds a message to the block.
    AddMessage(Message)

    // Digest computes the hash of the block.
    Digest() []byte

    // Valid returns true if the block is valid.
    Valid() bool
}
