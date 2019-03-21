package lazyledger

import (
    "bytes"
    "crypto/sha256"

    "gitlab.com/NebulousLabs/merkletree"
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

// PrevHash returns the hash of the previous block.
func (sb *SimpleBlock) PrevHash() []byte {
    return sb.prevHash
}

// Messages returns the block's messages.
func (sb *SimpleBlock) Messages() []Message {
    return sb.messages
}

// ApplicationProof creates a Merkle proof for all of the messages in a block for an application namespace.
// TODO: Deal with case to prove that there is no relevant messages in the block.
func (sb *SimpleBlock) ApplicationProof(namespace [namespaceSize]byte) (int, int, [][]byte) {
    var proofStart int
    var proofEnd int
    var found bool
    for index, message := range sb.messages {
        if message.Namespace() == namespace {
            if !found {
                found = true
                proofStart = index
            }
            proofEnd = index
        }
    }

    ndf := NewNamespaceDummyFlagger()
    fh := NewFlagHasher(ndf, sha256.New())
    proof, _ := merkletree.BuildRangeProof(proofStart, proofEnd, NewMessageSubtreeHasher(&sb.messages, fh))
    return proofStart, proofEnd, proof
}
