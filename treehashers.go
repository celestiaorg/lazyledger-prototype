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

// NewMessageSubtreeHasher returns a new NewMessageSubtreeHasher for a pointer to a slice of messages.
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
    return
}

// sum returns the hash of the input data using the specified algorithm.
func sum(h hash.Hash, data ...[]byte) []byte {
	h.Reset()
	for _, d := range data {
		// the Hash interface specifies that Write never returns an error
		_, _ = h.Write(d)
	}
	return h.Sum(nil)
}

// leafSum returns the hash created from data inserted to form a leaf. Leaf
// sums are calculated using:
//		Hash(0x00 || data)
func leafSum(h hash.Hash, data []byte) []byte {
	return sum(h, []byte{0x00}, data)
}
