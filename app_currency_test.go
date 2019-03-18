package lazyledger

import (
    "crypto/rand"
    "encoding/binary"
    "testing"

    "github.com/libp2p/go-libp2p-crypto"
)

func TestAppCurrency(t *testing.T) {
    bs := NewSimpleBlockStore()
    b := NewBlockchain(bs)

    sb := NewSimpleBlock([]byte{0})

    ms := NewSimpleMap()
    app := NewCurrency(ms)
    var regApp Application
    regApp = app
    b.RegisterApplication(&regApp)

    privA, pubA, _ := crypto.GenerateSecp256k1Key(rand.Reader)
    _, pubB, _ := crypto.GenerateSecp256k1Key(rand.Reader)
    pubABytes, _ := pubA.Bytes()
    pubABalanceBytes := make([]byte, binary.MaxVarintLen64)
    binary.BigEndian.PutUint64(pubABalanceBytes, 1000)
    ms.Put(pubABytes, pubABalanceBytes)

    sb.AddMessage(app.GenerateTransaction(privA, pubB, 100))
    b.ProcessBlock(sb)

    if app.Balance(pubA) != 900 || app.Balance(pubB) != 100 {
        t.Errorf("test tranasaction failed: invalid post-balances")
    }
}
