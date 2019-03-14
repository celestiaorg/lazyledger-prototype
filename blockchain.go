package lazyledger

// Blockchain is a chain of blocks.
type Blockchain struct {
    blockStore BlockStore
}

// NewBlockchain returns a new blockchain.
func NewBlockchain(blockStore BlockStore) *Blockchain {
    return &Blockchain{
        blockStore: blockStore,
    }
}

// PutBlock adds a new block to the block store.
func (b *Blockchain) PutBlock(block Block) {
    b.blockStore.Put(block.Digest(), block)
}
