package lazyledger

import (
    "encoding/binary"

    "github.com/golang/protobuf/proto"
    "github.com/libp2p/go-libp2p-crypto"
)

// Currency is a demo cryptocurrency application.
type Currency struct {
    state MapStore
    b *Blockchain
    transferCallbacks []TransferCallback
}

type TransferCallback = func(from []byte, to []byte, value int)

// NewCurrency creates a new currency instance.
func NewCurrency(state MapStore, b *Blockchain) *Currency {
    return &Currency{
        state: state,
        b: b,
    }
}

// ProcessMessage processes a message.
func (c *Currency) ProcessMessage(message Message) {
    transaction := &CurrencyTransaction{}
    err := proto.Unmarshal(message.Data(), transaction)
    if err != nil {
        return
    }
    if transaction.Dependency != nil {
        block, err := c.b.Block(c.BlockHead())
        if err != nil {
            return
        }
        dependencyProven := block.DependencyProven(transaction.Dependency)
        if !dependencyProven {
            return
        }
    }
    transactionMessage := &CurrencyTransactionMessage{
        To: transaction.To,
        Amount: transaction.Amount,
        Dependency: transaction.Dependency,
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

// GenerateTransaction generates a transaction message.
func (c *Currency) GenerateTransaction(fromPrivKey crypto.PrivKey, toPubKey crypto.PubKey, amount uint64, dependency []byte) Message {
    toPubKeyBytes, _ := toPubKey.Bytes()
    fromPubKeyBytes, _ := fromPrivKey.GetPublic().Bytes()
    transactionMessage := &CurrencyTransactionMessage{
        To: toPubKeyBytes,
        Amount: &amount,
        Dependency: dependency,
    }
    signedData, _ := proto.Marshal(transactionMessage)
    signature, _ := fromPrivKey.Sign(signedData)
    transaction := &CurrencyTransaction{
        To: toPubKeyBytes,
        From: fromPubKeyBytes,
        Amount: &amount,
        Signature: signature,
        Dependency: dependency,
    }
    messageData, _ := proto.Marshal(transaction)
    return *NewMessage(c.Namespace(), messageData)
}

// Balance gets the balance of a public key.
func (c *Currency) Balance(pubKey crypto.PubKey) uint64 {
    pubKeyBytes, _ := pubKey.Bytes()
    balance, err := c.state.Get(pubKeyBytes)
    if err != nil {
        return 0
    }
    return binary.BigEndian.Uint64(balance)
}

func (c *Currency) AddTransfer(fn TransferCallback) {
    c.transferCallbacks = append(c.transferCallbacks, fn)
}

func (c *Currency) triggerTransferCallbacks(from []byte, to []byte, value int) {
    for _, fn := range c.transferCallbacks {
        fn(from, to, value)
    }
}
