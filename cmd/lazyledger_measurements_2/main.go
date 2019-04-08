package main

import (
    "crypto/rand"
    "encoding/binary"
    "fmt"

    "github.com/musalbas/lazyledger-prototype"
    "github.com/libp2p/go-libp2p-crypto"
)

const namespaceSize = 8

func main() {
    currencyTxes := 10
    txAmounts := []int{128, 2048}
    txSize := 256
    for _, txes := range txAmounts {
        sb, ns := generateSimpleBlock(currencyTxes, txes, txSize)
        _, _, proofs1, messages1, _ := sb.ApplicationProof(ns)
        sbBandwidth := 0
        for _, msg := range *messages1 {
            sbBandwidth += len(msg.Marshal())
        }
        for _, hash := range proofs1 {
            sbBandwidth += len(hash)
        }

        pb, ns := generateProbabilisticBlock(currencyTxes, txes, txSize)
        _, _, proofs2, messages2, _ := pb.ApplicationProof(ns)
        pbBandwidth := 0
        for _, msg := range *messages2 {
            pbBandwidth += len(msg.Marshal())
        }
        for _, proof := range proofs2 {
            for _, hash := range proof {
                pbBandwidth += len(hash)
            }
        }

        fmt.Println(txes, sbBandwidth, pbBandwidth)
    }
}

func generateSimpleBlock(currencyTxes int, otherTxes, txSize int) (*lazyledger.SimpleBlock, [namespaceSize]byte) {
    txSize -= namespaceSize

    bs := lazyledger.NewSimpleBlockStore()
    b := lazyledger.NewBlockchain(bs)
    sb := lazyledger.NewSimpleBlock([]byte{0})
    ms := lazyledger.NewSimpleMap()
    app := lazyledger.NewCurrency(ms, b)
    b.RegisterApplication(&app)

    privA, pubA, _ := crypto.GenerateSecp256k1Key(rand.Reader)
    _, pubB, _ := crypto.GenerateSecp256k1Key(rand.Reader)
    pubABytes, _ := pubA.Bytes()
    pubABalanceBytes := make([]byte, binary.MaxVarintLen64)
    binary.BigEndian.PutUint64(pubABalanceBytes, 1000000)
    ms.Put(pubABytes, pubABalanceBytes)

    for i := 0; i < currencyTxes; i++ {
        sb.AddMessage(app.(*lazyledger.Currency).GenerateTransaction(privA, pubB, 1, nil))
    }

    for i := 0; i < otherTxes; i++ {
        messageData := make([]byte, txSize)
        rand.Read(messageData)
        sb.AddMessage(*lazyledger.NewMessage([namespaceSize]byte{0}, messageData))
    }

    return sb.(*lazyledger.SimpleBlock), app.Namespace()
}

func generateProbabilisticBlock(currencyTxes int, otherTxes, txSize int) (*lazyledger.ProbabilisticBlock, [namespaceSize]byte) {
    txSize -= namespaceSize

    bs := lazyledger.NewSimpleBlockStore()
    b := lazyledger.NewBlockchain(bs)
    pb := lazyledger.NewProbabilisticBlock([]byte{0}, 256)
    ms := lazyledger.NewSimpleMap()
    app := lazyledger.NewCurrency(ms, b)
    b.RegisterApplication(&app)

    privA, pubA, _ := crypto.GenerateSecp256k1Key(rand.Reader)
    _, pubB, _ := crypto.GenerateSecp256k1Key(rand.Reader)
    pubABytes, _ := pubA.Bytes()
    pubABalanceBytes := make([]byte, binary.MaxVarintLen64)
    binary.BigEndian.PutUint64(pubABalanceBytes, 1000000)
    ms.Put(pubABytes, pubABalanceBytes)

    for i := 0; i < currencyTxes; i++ {
        pb.AddMessage(app.(*lazyledger.Currency).GenerateTransaction(privA, pubB, 1, nil))
    }

    for i := 0; i < otherTxes; i++ {
        messageData := make([]byte, txSize)
        rand.Read(messageData)
        pb.AddMessage(*lazyledger.NewMessage([namespaceSize]byte{0}, messageData))
    }

    return pb.(*lazyledger.ProbabilisticBlock), app.Namespace()
}
