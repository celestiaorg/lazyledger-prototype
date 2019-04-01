package lazyledger

import (
    "hash"
    "io"

    "gitlab.com/NebulousLabs/merkletree"
)

// MessageSubtreeHasher implements merkletree.SubtreeHasher by reading leaf data from an underlying slice of messages.
type MessageSubtreeHasher struct {
    messages *[]Message
    h hash.Hash
    pos int
}

// NewMessageSubtreeHasher returns a new MessageSubtreeHasher for a pointer to a slice of messages.
func NewMessageSubtreeHasher(messages *[]Message, h hash.Hash) *MessageSubtreeHasher {
    return &MessageSubtreeHasher{
        messages: messages,
        h: h,
    }
}

// NextSubtreeRoot implements SubtreeHasher.
func (msh *MessageSubtreeHasher) NextSubtreeRoot(subtreeSize int) ([]byte, error) {
    tree := merkletree.New(msh.h)
    for i := 0; i < subtreeSize; i++ {
        if msh.pos >= len(*msh.messages) {
            break
        }
        tree.Push((*msh.messages)[msh.pos].Marshal())
        msh.pos += 1
    }

    root := tree.Root()
    if root == nil {
        return nil, io.EOF
    }

    return root, nil
}

// Skip implements SubtreeHasher.
func (msh *MessageSubtreeHasher) Skip(n int) error {
    if msh.pos + n > len(*msh.messages) {
        msh.pos = len(*msh.messages)
        return io.ErrUnexpectedEOF
    }
    msh.pos += n
    return nil
}

// MessageLeafHasher implements the merkletree.LeafHasher interface by reading leaf data from an underlying slice of messages.
type MessageLeafHasher struct {
    messages *[]Message
    h hash.Hash
    pos int
}

// NewMessageLeafHasher returns a new MessageLeafHasher for a pointer to a slice of messages.
func NewMessageLeafHasher(messages *[]Message, h hash.Hash) *MessageLeafHasher {
    return &MessageLeafHasher{
        messages: messages,
        h: h,
    }
}

// NextLeafHash implements LeafHasher
func (mlh *MessageLeafHasher) NextLeafHash() (leafHash []byte, err error) {
    if mlh.pos >= len(*mlh.messages) {
        return nil, io.EOF
    }

    leafHash = leafSum(mlh.h, (*mlh.messages)[mlh.pos].Marshal())
    err = nil
    mlh.pos += 1
    return
}

// HashLeafHasher implements the merkletree.LeafHasher interface from a slice of hashes.
type HashLeafHasher struct {
    hashes [][]byte
    pos int
}

// NewHashLeafHasher returns a new MessageLeafHasher for a slice of hashes.
func NewHashLeafHasher(hashes [][]byte) *HashLeafHasher {
    return &HashLeafHasher{
        hashes: hashes,
    }
}

// NextLeafHash implements LeafHasher
func (hlh *HashLeafHasher) NextLeafHash() (leafHash []byte, err error) {
    if hlh.pos >= len(hlh.hashes) {
        return nil, io.EOF
    }

    leafHash = hlh.hashes[hlh.pos]
    err = nil
    hlh.pos += 1
    return
}

// BytesLeafHasher implements the merkletree.LeafHasher interface by reading a slice.
type BytesLeafHasher struct {
    data [][]byte
    h hash.Hash
    pos int
}

// NewBytesLeafHasher returns a new BytesLeafHasher for a slice.
func NewBytesLeafHasher(data [][]byte, h hash.Hash) *BytesLeafHasher {
    return &BytesLeafHasher{
        data: data,
        h: h,
    }
}

// NextLeafHash implements LeafHasher
func (blh *BytesLeafHasher) NextLeafHash() (leafHash []byte, err error) {
    if blh.pos >= len(blh.data) {
        return nil, io.EOF
    }

    leafHash = leafSum(blh.h, blh.data[blh.pos])
    err = nil
    blh.pos += 1
    return
}

//CodedAxisSubtreeHasher implements merkletree.SubtreeHasher by reading leaf data from an underlying slice representing a coded axis.
type CodedAxisSubtreeHasher struct {
    data [][]byte
    h hash.Hash
    pos int
}

// NewCodedAxisSubtreeHasher returns a new CodedAxisSubtreeHasher for a pointer to a slice representing a coded axis.
func NewCodedAxisSubtreeHasher(data [][]byte, h hash.Hash) *CodedAxisSubtreeHasher {
    return &CodedAxisSubtreeHasher{
        data: data,
        h: h,
    }
}

// NextSubtreeRoot implements SubtreeHasher.
func (cash *CodedAxisSubtreeHasher) NextSubtreeRoot(subtreeSize int) ([]byte, error) {
    tree := merkletree.New(cash.h)
    for i := 0; i < subtreeSize; i++ {
        if cash.pos >= len(cash.data) {
            break
        }
        if cash.pos >= len(cash.data) / 2 {
            cash.h.(*flagDigest).setCodedMode(true)
        }
        tree.Push(cash.data[cash.pos])
        cash.h.(*flagDigest).setCodedMode(false)
        cash.pos += 1
    }

    root := tree.Root()
    if root == nil {
        return nil, io.EOF
    }

    return root, nil
}

// Skip implements SubtreeHasher.
func (cash *CodedAxisSubtreeHasher) Skip(n int) error {
    if cash.pos + n > len(cash.data) {
        cash.pos = len(cash.data)
        return io.ErrUnexpectedEOF
    }
    cash.pos += n
    return nil
}
