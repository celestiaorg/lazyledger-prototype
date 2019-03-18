package lazyledger

import (
    "encoding/binary"

    "github.com/golang/protobuf/proto"
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
func (c *Currency) ProcessMessage(message Message) {
    transaction := &CurrencyTransaction{}
    err := proto.Unmarshal(message.Data(), transaction)
    if err != nil {
        return
    }
    transactionMessage := &CurrencyTransactionMessage{
        To: transaction.To,
        Amount: transaction.Amount,
    }
    signedData, err := proto.Marshal(transactionMessage)
    fromKey, err := crypto.UnmarshalPublicKey(transaction.From)
    ok, err := fromKey.Verify(signedData, transaction.Signature)
    fromBalanceBytes, err := c.state.Get(transaction.From)
    if err != nil {
        return
    }
    fromBalance := binary.BigEndian.Uint64(fromBalanceBytes)
    if ok && fromBalance >= *transaction.Amount {
        toBalanceBytes, err := c.state.Get(transaction.To)
        var toBalance uint64
        if err != nil {
            toBalance = 0
        } else {
            toBalance = binary.BigEndian.Uint64(toBalanceBytes)
        }
        newFromBalanceBytes := make([]byte, binary.MaxVarintLen64)
        binary.BigEndian.PutUint64(newFromBalanceBytes, fromBalance - *transaction.Amount)
        newToBalanceBytes := make([]byte, binary.MaxVarintLen64)
        binary.BigEndian.PutUint64(newToBalanceBytes, toBalance + *transaction.Amount)

        c.state.Put(transaction.From, newFromBalanceBytes)
        c.state.Put(transaction.To, newToBalanceBytes)
    }
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
    head, _ := c.state.Get([]byte("__head__"))
    return head
}
