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
    provenDependencies map[string]bool
}

// NewSimpleBlock returns a new simple block.
func NewSimpleBlock(prevHash []byte) Block {
    return &SimpleBlock{
        prevHash: prevHash,
        provenDependencies: make(map[string]bool),
    }
}

// ImportSimpleBlockHeader imports a received simple block without the messages.
func ImportSimpleBlockHeader(prevHash []byte, messagesRoot []byte) Block {
    return &SimpleBlock{
        prevHash: prevHash,
        messagesRoot: messagesRoot,
        provenDependencies: make(map[string]bool),
    }
}

// ImportSimpleBlock imports a received simple block.
func ImportSimpleBlock(prevHash []byte, messages []Message) Block {
    return &SimpleBlock{
        prevHash: prevHash,
        messages: messages,
        provenDependencies: make(map[string]bool),
    }
}

// AddMessage adds a message to the block.
func (sb *SimpleBlock) AddMessage(message Message) {
    sb.messages = append(sb.messages, message)

    // Force recompututation of messagesRoot
    sb.messagesRoot = nil
}

// MessagesRoot returns the Merkle root of the messages in the block.
func (sb *SimpleBlock) MessagesRoot() []byte {
    if sb.messagesRoot == nil {
        ndf := NewNamespaceDummyFlagger()
        fh := NewFlagHasher(ndf, sha256.New())
        tree := merkletree.New(fh)
        for _, message := range sb.messages {
            tree.Push(message.Marshal())
        }
        sb.messagesRoot = tree.Root()
    }

    return sb.messagesRoot
}

// Digest computes the hash of the block.
func (sb *SimpleBlock) Digest() []byte {
    hasher := sha256.New()
    hasher.Write(sb.prevHash)
    hasher.Write(sb.MessagesRoot())
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
    if bytes.Compare(tree.Root(), sb.MessagesRoot()) == 0 {
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
func (sb *SimpleBlock) ApplicationProof(namespace [namespaceSize]byte) (int, int, [][]byte, *[]Message, [][]byte) {
    var proofStart int
    var proofEnd int
    var found bool
    for index, message := range sb.messages {
        if message.Namespace() == namespace {
            if !found {
                found = true
                proofStart = index
            }
            proofEnd = index + 1
        }
    }

    var inRange bool
    if !found {
        var prevMessage Message
        // We need to generate a proof for an absence of relevant messages.
        for index, message := range sb.messages {
            if index != 0 {
                prevNs := prevMessage.Namespace()
                currentNs := message.Namespace()
                if ((bytes.Compare(prevNs[:], namespace[:]) < 0 && bytes.Compare(namespace[:], currentNs[:]) < 0) ||
                    (bytes.Compare(prevNs[:], namespace[:]) > 0 && bytes.Compare(namespace[:], currentNs[:]) > 0)) {
                    if !inRange {
                        inRange = true
                        proofStart = index
                    }
                    proofEnd = index + 1
                }
            }
            prevMessage = message
        }
    }

    ndf := NewNamespaceDummyFlagger()
    fh := NewFlagHasher(ndf, sha256.New())
    var proof [][]byte
    if found || inRange {
        proof, _ = merkletree.BuildRangeProof(proofStart, proofEnd, NewMessageSubtreeHasher(&sb.messages, fh))
    }
    proofMessages := sb.messages[proofStart:proofEnd]
    if found {
        return proofStart, proofEnd, proof, &proofMessages, nil
    }

    var hashes [][]byte
    for _, message := range proofMessages {
        ndf := NewNamespaceDummyFlagger()
        fh := NewFlagHasher(ndf, sha256.New())
        hashes = append(hashes, leafSum(fh, message.Marshal()))
        fh.Reset()
    }

    return proofStart, proofEnd, proof, nil, hashes
}

// VerifyApplicationProof verifies a Merkle proof for all of the messages in a block for an application namespace.
func (sb *SimpleBlock) VerifyApplicationProof(namespace [namespaceSize]byte, proofStart int, proofEnd int, proof [][]byte, messages *[]Message, hashes [][]byte) bool {
    // Verify Merkle proof
    ndf := NewNamespaceDummyFlagger()
    fh := NewFlagHasher(ndf, sha256.New())
    var lh merkletree.LeafHasher
    if messages != nil {
        lh = NewMessageLeafHasher(messages, fh)
    } else {
        lh = NewHashLeafHasher(hashes)
    }
    result, err := merkletree.VerifyRangeProof(lh, fh, proofStart, proofEnd, proof, sb.MessagesRoot())
    if !result || err != nil {
        return false
    }

    // Verify proof completeness
    var leafIndex uint64
    var leftSubtrees [][]byte
    var rightSubtrees [][]byte
	consumeUntil := func(end uint64) error {
		for leafIndex != end && len(proof) > 0 {
			subtreeSize := nextSubtreeSize(leafIndex, end)
            leftSubtrees = append(leftSubtrees, proof[0])
			proof = proof[1:]
			leafIndex += uint64(subtreeSize)
		}
		return nil
	}
    if err := consumeUntil(uint64(proofStart)); err != nil {
        return false
    }
    rightSubtrees = proof

    for _, subtree := range leftSubtrees {
        _, max := dummyNamespacesFromFlag(subtree)
        if bytes.Compare(max, namespace[:]) >= 0 {
            return false
        }
    }
    for _, subtree := range rightSubtrees {
        min, _ := dummyNamespacesFromFlag(subtree)
        if bytes.Compare(min, namespace[:]) <= 0 {
            return false
        }
    }

    return true
}

func (sb *SimpleBlock) ProveDependency(index int) ([]byte, [][]byte, error) {
    ndf := NewNamespaceDummyFlagger()
    fh := NewFlagHasher(ndf, sha256.New())
    proof, err := merkletree.BuildRangeProof(index, index + 1, NewMessageSubtreeHasher(&sb.messages, fh))
    if err != nil {
        return nil, nil, err
    }
    return leafSum(fh, sb.messages[index].Marshal()), proof, nil
}

func (sb *SimpleBlock) VerifyDependency(index int, hash []byte, proof [][]byte) bool {
    ndf := NewNamespaceDummyFlagger()
    fh := NewFlagHasher(ndf, sha256.New())
    lh := NewHashLeafHasher([][]byte{hash})
    result, err := merkletree.VerifyRangeProof(lh, fh, index, index + 1, proof, sb.MessagesRoot())
    if result && err == nil {
        sb.provenDependencies[string(hash)] = true
        return true
    }
    return false
}

func (sb *SimpleBlock) DependencyProven(hash []byte) bool {
    if value, ok := sb.provenDependencies[string(hash)]; ok {
        return value
    }
    return false
}
