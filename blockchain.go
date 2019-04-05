package lazyledger

import (
    "bytes"
)

// Blockchain is a chain of blocks.
// This is a prototype for testing purposes and thus does not support re-orgs, and there is no network stack.
type Blockchain struct {
    blockStore BlockStore
    headBlock Block
    applications []*Application
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
        b.processCallbacks(block, true)
    } else {
        b.processCallbacks(block, false)
    }
}

func (b *Blockchain) Block(digest []byte) (Block, error) {
    block, err := b.blockStore.Get(digest)
    if err != nil {
        return nil, err
    }
    return block, nil
}

// RegisterApplication registers an application instance to call when new relevant messages arrive.
func (b *Blockchain) RegisterApplication(application *Application) {
    b.applications = append(b.applications, application)
}

func (b *Blockchain) processCallbacks(block Block, isHead bool) {
    for _, application := range b.applications {
        if isHead {
            (*application).SetBlockHead(block.Digest())
        }
        for _, message := range block.Messages() {
            if message.Namespace() == (*application).Namespace() {
                (*application).ProcessMessage(message)
            }
        }
    }
}
