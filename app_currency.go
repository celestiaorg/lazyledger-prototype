package lazyledger

import (
    "github.com/libp2p/go-libp2p-crypto"
)

// Demo cryptocurrency application.
type Currency struct {
    state MapStore
}

// NewCurrency creates a new currency instance.
func NewCurrency(state MapStore) *Currency {
    return &Currency{
        state: state,
    }
}

// ProcessMessage processes a message.
func (c *Currency) ProcessMessage(message Message) bool {
    return true
}

// Namespace returns the application's namespace ID.
func (c *Currency) Namespace() [namespaceSize]byte {
    var namespace [namespaceSize]byte
    copy(namespace[:], []byte("currency"))
    return namespace
}

// SetBlockHead sets the hash of the latest block that has been processed.
func (c *Currency) SetBlockHead(hash []byte) {
    c.state.Put([]byte("__head__"), hash)
}

// BlockHead returns the hash of the latest block that has been processed.
func (c *Currency) BlockHead() []byte {
    return c.state.Get([]byte("__head__"))
}
