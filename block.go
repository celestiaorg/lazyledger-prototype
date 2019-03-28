package lazyledger

// Block represents a block in the chain.
type Block interface {
    // AddMessage adds a message to the block.
    AddMessage(Message)

    // Digest computes the hash of the block.
    Digest() []byte

    // Valid returns true if the block is valid.
    Valid() bool

    // PrevHash returns the hash of the previous block.
    PrevHash() []byte

    // Messages returns the block's messages.
    Messages() []Message

    // ApplicationProof creates a Merkle proof for all of the messages in a block for an application namespace.
    ApplicationProof([namespaceSize]byte) (int, int, [][]byte, *[]Message, [][]byte)

    // VerifyApplicationProof verifies a Merkle proof for all of the messages in a block for an application namespace.
    VerifyApplicationProof([namespaceSize]byte, int, int, [][]byte, *[]Message, [][]byte) bool
}
