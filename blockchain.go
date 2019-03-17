package lazyledger

import (
    "bytes"
)

// Blockchain is a chain of blocks.
// This is a prototype for testing purposes and thus does not support re-orgs, and there is no network stack.
type Blockchain struct {
    blockStore BlockStore
    headBlock Block
}

// NewBlockchain returns a new blockchain.
func NewBlockchain(blockStore BlockStore) *Blockchain {
    return &Blockchain{
        blockStore: blockStore,
    }
}

// ProcessBlock processes a new block.
func (b *Blockchain) ProcessBlock(block Block) {
    b.blockStore.Put(block.Digest(), block)

    if b.headBlock == nil || bytes.Compare(block.PrevHash(), b.headBlock.Digest()) == 0 {
        b.headBlock = block
    }
}
