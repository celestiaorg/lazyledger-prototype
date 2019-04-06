package lazyledger

import (
    "crypto/rand"
    "encoding/binary"
    "testing"
    "bytes"

    "github.com/libp2p/go-libp2p-crypto"
)

func TestAppRegistrarSimpleBlock(t *testing.T) {
    bs := NewSimpleBlockStore()
    b := NewBlockchain(bs)

    sb := NewSimpleBlock([]byte{0})

    ms1 := NewSimpleMap()
    currencyApp := NewCurrency(ms1, b)
    b.RegisterApplication(&currencyApp)

    privA, pubA, _ := crypto.GenerateSecp256k1Key(rand.Reader)
    _, pubB, _ := crypto.GenerateSecp256k1Key(rand.Reader)
    pubABytes, _ := pubA.Bytes()
    pubBBytes, _ := pubB.Bytes()
    pubABalanceBytes := make([]byte, binary.MaxVarintLen64)
    binary.BigEndian.PutUint64(pubABalanceBytes, 1000)
    ms1.Put(pubABytes, pubABalanceBytes)

    ms2 := NewSimpleMap()
    registrarApp := NewRegistrar(ms2, currencyApp.(*Currency), pubBBytes)
    b.RegisterApplication(&registrarApp)

    sb.AddMessage(currencyApp.(*Currency).GenerateTransaction(privA, pubB, 100, nil))
    sb.AddMessage(registrarApp.(*Registrar).GenerateTransaction(privA, []byte("foo")))
    b.ProcessBlock(sb)

    if currencyApp.(*Currency).Balance(pubA) != 900 || currencyApp.(*Currency).Balance(pubB) != 100 {
        t.Error("test tranasaction failed: invalid post-balances")
    }
    if registrarApp.(*Registrar).Balance(pubABytes) != 100 {
        t.Error("test tranasaction failed: invalid post-balances in registrar")
    }
    if bytes.Compare(registrarApp.(*Registrar).Name([]byte("foo")), pubABytes) != 0 {
        t.Error("failed to register name")
    }
}
