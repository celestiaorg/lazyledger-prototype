package lazyledger

import (
    "bytes"
    "crypto/sha256"

    "github.com/NebulousLabs/merkletree"
)

// SimpleBlock represents a block designed for the Simple Validity Rule.
type SimpleBlock struct {
    prevHash []byte
    messages []Message
    messagesRoot []byte
}

// NewSimpleBlock returns a new simple block.
func NewSimpleBlock(prevHash []byte) Block {
    return &SimpleBlock{
        prevHash: prevHash,
    }
}

// ImportSimpleBlockHeader imports a received simple block without the messages.
func ImportSimpleBlockHeader(prevHash []byte, messagesRoot []byte) Block {
    return &SimpleBlock{
        prevHash: prevHash,
        messagesRoot: messagesRoot,
    }
}

// ImportSimpleBlock imports a received simple block.
func ImportSimpleBlock(prevHash []byte, messages []Message) Block {
    return &SimpleBlock{
        prevHash: prevHash,
        messages: messages,
    }
}

// AddMessage adds a message to the block.
func (sb *SimpleBlock) AddMessage(message Message) {
    sb.messages = append(sb.messages, message)
}

// Digest computes the hash of the block.
func (sb *SimpleBlock) Digest() []byte {
    hasher := sha256.New()
    hasher.Write(sb.prevHash)
    hasher.Write(sb.messagesRoot)
    return hasher.Sum(nil)
}

// Valid returns true if the block is valid.
func (sb *SimpleBlock) Valid() bool {
    if sb.messages == nil {
        // Cannot validate block without messages.
        return false
    }

    ndf := NewNamespaceDummyFlagger()
    fh := NewFlagHasher(ndf, sha256.New())
    tree := merkletree.New(fh)
    for _, message := range sb.messages {
        tree.Push(message.Marshal())
    }
    if bytes.Compare(tree.Root(), sb.messagesRoot) == 0 {
        return true
    }
    return false
}

func (sb *SimpleBlock) PrevHash() []byte {
    return sb.prevHash
}
