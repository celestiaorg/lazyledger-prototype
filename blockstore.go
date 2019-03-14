package lazyledger

import(
    "fmt"
)

// BlockStore is a key-value store for blocks.
type BlockStore interface {
    Get(key []byte) (Block, error) // Get gets the value for a key.
    Put(key []byte, value Block) error // Put updates the value for a key.
    Del(key []byte) error // Del deletes a key.
}

// InvalidKeyError is thrown when a key that does not exist is being accessed.
type InvalidKeyError struct {
    Key []byte
}

func (e *InvalidKeyError) Error() string {
    return fmt.Sprintf("invalid key: %s", e.Key)
}

// SimpleBlockStore is a simple in-memory block store.
type SimpleBlockStore struct {
    m map[string]Block
}

// NewSimpleBlockStore creates a new empty simple block store.
func NewSimpleBlockStore() *SimpleBlockStore {
    return &SimpleBlockStore{
        m: make(map[string]Block),
    }
}

// Get gets the value for a key.
func (sbs *SimpleBlockStore) Get(key []byte) (Block, error) {
    if value, ok := sbs.m[string(key)]; ok {
        return value, nil
    }
    return nil, &InvalidKeyError{Key: key}
}

// Put updates the value for a key.
func (sbs *SimpleBlockStore) Put(key []byte, value Block) error {
    sbs.m[string(key)] = value
    return nil
}

// Del deletes a key.
func (sbs *SimpleBlockStore) Del(key []byte) error {
    _, ok := sbs.m[string(key)]
    if ok {
        delete(sbs.m, string(key))
        return nil
    }
    return &InvalidKeyError{Key: key}
}
