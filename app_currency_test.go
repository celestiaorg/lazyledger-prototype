package lazyledger

import (
    "crypto/rand"
    "encoding/binary"
    "testing"

    "github.com/libp2p/go-libp2p-crypto"
)

func TestAppCurrencySimpleBlock(t *testing.T) {
    bs := NewSimpleBlockStore()
    b := NewBlockchain(bs)

    sb := NewSimpleBlock([]byte{0})

    ms := NewSimpleMap()
    app := NewCurrency(ms, b)
    b.RegisterApplication(&app)

    privA, pubA, _ := crypto.GenerateSecp256k1Key(rand.Reader)
    _, pubB, _ := crypto.GenerateSecp256k1Key(rand.Reader)
    pubABytes, _ := pubA.Bytes()
    pubABalanceBytes := make([]byte, binary.MaxVarintLen64)
    binary.BigEndian.PutUint64(pubABalanceBytes, 1000)
    ms.Put(pubABytes, pubABalanceBytes)

    sb.AddMessage(app.(*Currency).GenerateTransaction(privA, pubB, 100, nil))
    b.ProcessBlock(sb)

    if app.(*Currency).Balance(pubA) != 900 || app.(*Currency).Balance(pubB) != 100 {
        t.Error("test tranasaction failed: invalid post-balances")
    }
}

func TestAppCurrencySimpleBlockDependency(t *testing.T) {
    bs := NewSimpleBlockStore()
    b := NewBlockchain(bs)

    sb := NewSimpleBlock([]byte{0})

    ms := NewSimpleMap()
    app := NewCurrency(ms, b)
    b.RegisterApplication(&app)

    privA, pubA, _ := crypto.GenerateSecp256k1Key(rand.Reader)
    _, pubB, _ := crypto.GenerateSecp256k1Key(rand.Reader)
    pubABytes, _ := pubA.Bytes()
    pubABalanceBytes := make([]byte, binary.MaxVarintLen64)
    binary.BigEndian.PutUint64(pubABalanceBytes, 1000)
    ms.Put(pubABytes, pubABalanceBytes)

    sb.AddMessage(*NewMessage([namespaceSize]byte{0}, []byte("foo")))
    hash, _, _ := sb.ProveDependency(0)
    sb.AddMessage(app.(*Currency).GenerateTransaction(privA, pubB, 100, hash))
    _, proof, _ := sb.ProveDependency(0)
    sb.VerifyDependency(0, hash, proof)
    b.ProcessBlock(sb)

    if app.(*Currency).Balance(pubA) != 900 || app.(*Currency).Balance(pubB) != 100 {
        t.Error("test tranasaction failed: invalid post-balances")
    }
}

func TestAppCurrencyProbabilisticBlockDependency(t *testing.T) {
    bs := NewSimpleBlockStore()
    b := NewBlockchain(bs)

    pb := NewProbabilisticBlock([]byte{0}, 512)

    ms := NewSimpleMap()
    app := NewCurrency(ms, b)
    b.RegisterApplication(&app)

    privA, pubA, _ := crypto.GenerateSecp256k1Key(rand.Reader)
    _, pubB, _ := crypto.GenerateSecp256k1Key(rand.Reader)
    pubABytes, _ := pubA.Bytes()
    pubABalanceBytes := make([]byte, binary.MaxVarintLen64)
    binary.BigEndian.PutUint64(pubABalanceBytes, 1000)
    ms.Put(pubABytes, pubABalanceBytes)

    pb.AddMessage(*NewMessage([namespaceSize]byte{0}, []byte("foo")))
    hash, _, _ := pb.ProveDependency(0)
    pb.AddMessage(app.(*Currency).GenerateTransaction(privA, pubB, 100, hash))
    _, proof, _ := pb.ProveDependency(0)
    pb.VerifyDependency(0, hash, proof)
    b.ProcessBlock(pb)

    if app.(*Currency).Balance(pubA) != 900 || app.(*Currency).Balance(pubB) != 100 {
        t.Error("test tranasaction failed: invalid post-balances")
    }
}

func TestAppCurrencyProbabilisticBlock(t *testing.T) {
    bs := NewSimpleBlockStore()
    b := NewBlockchain(bs)

    pb := NewProbabilisticBlock([]byte{0}, 512)

    ms := NewSimpleMap()
    app := NewCurrency(ms, b)
    b.RegisterApplication(&app)

    privA, pubA, _ := crypto.GenerateSecp256k1Key(rand.Reader)
    _, pubB, _ := crypto.GenerateSecp256k1Key(rand.Reader)
    pubABytes, _ := pubA.Bytes()
    pubABalanceBytes := make([]byte, binary.MaxVarintLen64)
    binary.BigEndian.PutUint64(pubABalanceBytes, 1000)
    ms.Put(pubABytes, pubABalanceBytes)

    pb.AddMessage(app.(*Currency).GenerateTransaction(privA, pubB, 100, nil))
    b.ProcessBlock(pb)

    if app.(*Currency).Balance(pubA) != 900 || app.(*Currency).Balance(pubB) != 100 {
        t.Error("test tranasaction failed: invalid post-balances")
    }
}
