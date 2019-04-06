package lazyledger

import (
    "encoding/binary"
    "bytes"

    "github.com/golang/protobuf/proto"
    "github.com/libp2p/go-libp2p-crypto"
)

type Registrar struct {
    state MapStore
    currency *Currency
    owner []byte
}

func NewRegistrar(state MapStore, currency *Currency, owner []byte) Application {
    app := &Registrar{
        state: state,
        currency: currency,
        owner: owner,
    }
    currency.AddTransferCallback(app.transferCallback)
    return app
}

func (app *Registrar) ProcessMessage(message Message) {
    transaction := &RegisterTransaction{}
    err := proto.Unmarshal(message.Data(), transaction)
    if err != nil {
        return
    }
    transactionMessage := &RegisterTransactionMessage{
        Name: transaction.Name,
    }
    signedData, err := proto.Marshal(transactionMessage)
    ownerKey, err := crypto.UnmarshalPublicKey(transaction.Owner)
    ok, err := ownerKey.Verify(signedData, transaction.Signature)
    if !ok {
        return
    }
    if bytes.Compare(app.Name(transaction.Name), []byte{}) != 0 { // check name is available
        return
    }

    // check and subtract balance
    balance := app.Balance(transaction.Owner)
    if balance < 5 {
        return
    }
    newBalanceBytes := make([]byte, binary.MaxVarintLen64)
    binary.BigEndian.PutUint64(newBalanceBytes, balance - 5)
    app.state.Put(transaction.Owner, append([]byte("balance__"), newBalanceBytes...))

    app.state.Put(append([]byte("name__"), transaction.Name...), transaction.Owner)
}

func (app *Registrar) Namespace() [namespaceSize]byte {
    var namespace [namespaceSize]byte
    copy(namespace[:], []byte("reggie"))
    return namespace
}

func (app *Registrar) SetBlockHead(hash []byte) {
    app.state.Put([]byte("__head__"), hash)
}

func (app *Registrar) BlockHead() []byte {
    head, _ := app.state.Get([]byte("__head__"))
    return head
}

func (app *Registrar) Name(name []byte) []byte {
    value, err := app.state.Get(append([]byte("name__"), name...))
    if err != nil {
        return []byte{}
    }
    return value
}

func (app *Registrar) Balance(key []byte) uint64 {
    balance, err := app.state.Get(append([]byte("balance__"), key...))
    if err != nil {
        return 0
    }
    return binary.BigEndian.Uint64(balance)
}

func (app *Registrar) transferCallback(from []byte, to []byte, value int) {
    if bytes.Compare(app.owner, to) == 0 {
        balanceBytes, err := app.state.Get(append([]byte("balance__"), from...))
        var balance uint64
        if err != nil {
            balance = 0
        } else {
            balance = binary.BigEndian.Uint64(balanceBytes)
        }
        newBalanceBytes := make([]byte, binary.MaxVarintLen64)
        binary.BigEndian.PutUint64(newBalanceBytes, balance + uint64(value))
        app.state.Put(append([]byte("balance__"), from...), newBalanceBytes)
    }
}

func (app *Registrar) GenerateTransaction(owner crypto.PrivKey, name []byte) Message {
    ownerPubKeyBytes, _ := owner.GetPublic().Bytes()
    transactionMessage := &RegisterTransactionMessage{
        Name: name,
    }
    signedData, _ := proto.Marshal(transactionMessage)
    signature, _ := owner.Sign(signedData)
    transaction := &RegisterTransaction{
        Owner: ownerPubKeyBytes,
        Name: name,
        Signature: signature,
    }
    messageData, _ := proto.Marshal(transaction)
    return *NewMessage(app.Namespace(), messageData)
}
