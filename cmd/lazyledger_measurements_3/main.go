package main

import (
    "crypto/rand"
    "encoding/binary"
    "fmt"

    "github.com/LazyLedger/lazyledger-prototype"
    "github.com/libp2p/go-libp2p-crypto"
)

const namespaceSize = 8

func main() {
    currencyTxes := 10
    txAmounts := []int{128, 256, 384, 512, 640, 768, 896, 1024}
    txSize := 1024
    for _, txes := range txAmounts {
        _, c, d := generateSimpleBlock(currencyTxes, txes, txSize)
        sbStoragec := c.StorageSize()
        sbStoraged := d.StorageSize()

        //_, c, d = generateProbabilisticBlock(currencyTxes, txes, txSize)
        //pbStoragec := c.StorageSize()
        //pbStoraged := d.StorageSize()

        fmt.Println(txes, sbStoragec, sbStoraged)//, pbStoragec, pbStoraged)
    }
}

func generateSimpleBlock(currencyTxes int, otherTxes, txSize int) (*lazyledger.SimpleBlock, *lazyledger.Currency, *lazyledger.DummyApp) {
    bs := lazyledger.NewSimpleBlockStore()
    b := lazyledger.NewBlockchain(bs)
    sb := lazyledger.NewSimpleBlock([]byte{0})
    ms := lazyledger.NewSimpleMap()
    app := lazyledger.NewCurrency(ms, b)
    b.RegisterApplication(&app)

    ms2 := lazyledger.NewSimpleMap()
    app2 := lazyledger.NewDummyApp(ms2)
    b.RegisterApplication(&app2)

    privA, pubA, _ := crypto.GenerateSecp256k1Key(rand.Reader)
    pubABytes, _ := pubA.Bytes()
    pubABalanceBytes := make([]byte, binary.MaxVarintLen64)
    binary.BigEndian.PutUint64(pubABalanceBytes, 1000000)
    ms.Put(pubABytes, pubABalanceBytes)

    for i := 0; i < currencyTxes; i++ {
        _, pubB, _ := crypto.GenerateSecp256k1Key(rand.Reader)
        sb.AddMessage(app.(*lazyledger.Currency).GenerateTransaction(privA, pubB, 1, nil))
    }

    for i := 0; i < otherTxes; i++ {
        k := make([]byte, txSize / 2)
        v := make([]byte, txSize / 2)
        puts := make(map[string]string)
        rand.Read(k)
        rand.Read(v)
        puts[string(k)] = string(v)
        t := app2.(*lazyledger.DummyApp).GenerateTransaction(puts)
        sb.AddMessage(t)
    }

    b.ProcessBlock(sb)

    return sb.(*lazyledger.SimpleBlock), app.(*lazyledger.Currency), app2.(*lazyledger.DummyApp)
}

func generateProbabilisticBlock(currencyTxes int, otherTxes, txSize int) (*lazyledger.ProbabilisticBlock, *lazyledger.Currency, *lazyledger.DummyApp) {
    pb := lazyledger.NewProbabilisticBlock([]byte{0}, txSize+20)

    bs := lazyledger.NewSimpleBlockStore()
    b := lazyledger.NewBlockchain(bs)
    ms := lazyledger.NewSimpleMap()
    app := lazyledger.NewCurrency(ms, b)
    b.RegisterApplication(&app)

    privA, pubA, _ := crypto.GenerateSecp256k1Key(rand.Reader)
    pubABytes, _ := pubA.Bytes()
    pubABalanceBytes := make([]byte, binary.MaxVarintLen64)
    binary.BigEndian.PutUint64(pubABalanceBytes, 1000000)
    ms.Put(pubABytes, pubABalanceBytes)

    ms2 := lazyledger.NewSimpleMap()
    app2 := lazyledger.NewDummyApp(ms2)
    b.RegisterApplication(&app2)

    for i := 0; i < currencyTxes; i++ {
        _, pubB, _ := crypto.GenerateSecp256k1Key(rand.Reader)
        pb.AddMessage(app.(*lazyledger.Currency).GenerateTransaction(privA, pubB, 1, nil))
    }

    for i := 0; i < otherTxes; i++ {
        k := make([]byte, txSize / 2)
        v := make([]byte, txSize / 2)
        puts := make(map[string]string)
        rand.Read(k)
        rand.Read(v)
        puts[string(k)] = string(v)
        t := app2.(*lazyledger.DummyApp).GenerateTransaction(puts)
        pb.AddMessage(t)
    }

    b.ProcessBlock(pb)

    return pb.(*lazyledger.ProbabilisticBlock), app.(*lazyledger.Currency), app2.(*lazyledger.DummyApp)
}
