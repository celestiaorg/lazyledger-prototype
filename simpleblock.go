package lazyledger

import (
    "crypto/sha256"
)

// SimpleBlock represents a block designed for the Simple Validity Rule.
type SimpleBlock struct {
    prevHash []byte
    messages []Message
}

// NewSimpleBlock returns a new simple block.
func NewSimpleBlock(prevHash []byte) Block {
    return &SimpleBlock{
        prevHash: prevHash,
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
    for _, message := range sb.messages {
        hasher.Write(message.Marshal())
    }
    return hasher.Sum(nil)
}
