package main

import (
    "crypto/rand"
    "encoding/binary"
    "fmt"

    "github.com/lazyledger/lazyledger-prototype"
    "github.com/libp2p/go-libp2p-crypto"
)

const namespaceSize = 8

func main() {
    registrarTxes := 10
    txAmounts := []int{128, 256, 384, 512, 640, 768, 896, 1024}
    txSize := 256
    for _, txes := range txAmounts {
        sbBandwidthMax := 0
        pbBandwidthMax := 0

        for i := 0; i < 10; i++ {
            sb, cns, rns := generateSimpleBlock(registrarTxes, txes, txSize)
            sbBandwidth := 0
            _, _, proofs1, messages1, _ := sb.ApplicationProof(cns)
            for _, msg := range *messages1 {
                sbBandwidth += len(msg.Marshal())
            }
            for _, hash := range proofs1 {
                sbBandwidth += len(hash)
            }
            if sbBandwidth > sbBandwidthMax {
                sbBandwidthMax = sbBandwidth
            }
            _, _, proofs1, _, hashes1 := sb.ApplicationProof(rns)
            for _, hash := range hashes1 {
                sbBandwidth += len(hash)
            }
            for _, hash := range proofs1 {
                sbBandwidth += len(hash)
            }
            if sbBandwidth > sbBandwidthMax {
                sbBandwidthMax = sbBandwidth
            }

            pb, cns, rns := generateProbabilisticBlock(registrarTxes, txes, txSize)
            _, _, proofs2, messages2, _ := pb.ApplicationProof(cns)
            pbBandwidth := 0
            for _, msg := range *messages2 {
                pbBandwidth += len(msg.Marshal())
            }
            for _, proof := range proofs2 {
                for _, hash := range proof {
                    pbBandwidth += len(hash)
                }
            }
            _, _, proofs2, _, hashes2 := pb.ApplicationProof(rns)
            for _, hash := range hashes2 {
                pbBandwidth += len(hash)
            }
            for _, proof := range proofs2 {
                for _, hash := range proof {
                    pbBandwidth += len(hash)
                }
            }
            if pbBandwidth > pbBandwidthMax {
                pbBandwidthMax = pbBandwidth
            }
        }

        fmt.Println(txes, sbBandwidthMax, pbBandwidthMax)
    }
}

func generateSimpleBlock(registrarTxes int, otherTxes, txSize int) (*lazyledger.SimpleBlock, [namespaceSize]byte, [namespaceSize]byte) {
    txSize -= namespaceSize

    bs := lazyledger.NewSimpleBlockStore()
    b := lazyledger.NewBlockchain(bs)

    sb := lazyledger.NewSimpleBlock([]byte{0})

    ms1 := lazyledger.NewSimpleMap()
    currencyApp := lazyledger.NewCurrency(ms1, b)
    b.RegisterApplication(&currencyApp)

    privA, pubA, _ := crypto.GenerateSecp256k1Key(rand.Reader)
    _, pubB, _ := crypto.GenerateSecp256k1Key(rand.Reader)
    pubABytes, _ := pubA.Bytes()
    pubBBytes, _ := pubB.Bytes()
    pubABalanceBytes := make([]byte, binary.MaxVarintLen64)
    binary.BigEndian.PutUint64(pubABalanceBytes, 1000000)
    ms1.Put(pubABytes, pubABalanceBytes)

    ms2 := lazyledger.NewSimpleMap()
    registrarApp := lazyledger.NewRegistrar(ms2, currencyApp.(*lazyledger.Currency), pubBBytes)
    b.RegisterApplication(&registrarApp)

    for i := 0; i < registrarTxes; i++ {
        sb.AddMessage(currencyApp.(*lazyledger.Currency).GenerateTransaction(privA, pubB, 1, nil))
    }

    ms3 := lazyledger.NewSimpleMap()
    registrarApp2 := lazyledger.NewRegistrar(ms3, currencyApp.(*lazyledger.Currency), pubBBytes)
    b.RegisterApplication(&registrarApp2)

    for i := 0; i < otherTxes; i++ {
        sb.AddMessage(currencyApp.(*lazyledger.Currency).GenerateTransaction(privA, pubB, 1, nil))
    }

    return sb.(*lazyledger.SimpleBlock), currencyApp.Namespace(), registrarApp.Namespace()
}

func generateProbabilisticBlock(registrarTxes int, otherTxes, txSize int) (*lazyledger.ProbabilisticBlock, [namespaceSize]byte, [namespaceSize]byte) {
    pb := lazyledger.NewProbabilisticBlock([]byte{0}, txSize)
    txSize -= namespaceSize + 2

    bs := lazyledger.NewSimpleBlockStore()
    b := lazyledger.NewBlockchain(bs)

    ms1 := lazyledger.NewSimpleMap()
    currencyApp := lazyledger.NewCurrency(ms1, b)
    b.RegisterApplication(&currencyApp)

    privA, pubA, _ := crypto.GenerateSecp256k1Key(rand.Reader)
    _, pubB, _ := crypto.GenerateSecp256k1Key(rand.Reader)
    pubABytes, _ := pubA.Bytes()
    pubBBytes, _ := pubB.Bytes()
    pubABalanceBytes := make([]byte, binary.MaxVarintLen64)
    binary.BigEndian.PutUint64(pubABalanceBytes, 1000000)
    ms1.Put(pubABytes, pubABalanceBytes)

    ms2 := lazyledger.NewSimpleMap()
    registrarApp := lazyledger.NewRegistrar(ms2, currencyApp.(*lazyledger.Currency), pubBBytes)
    b.RegisterApplication(&registrarApp)

    for i := 0; i < registrarTxes; i++ {
        pb.AddMessage(currencyApp.(*lazyledger.Currency).GenerateTransaction(privA, pubB, 1, nil))
    }

    ms3 := lazyledger.NewSimpleMap()
    registrarApp2 := lazyledger.NewRegistrar(ms3, currencyApp.(*lazyledger.Currency), pubBBytes)
    b.RegisterApplication(&registrarApp2)

    for i := 0; i < otherTxes; i++ {
        pb.AddMessage(currencyApp.(*lazyledger.Currency).GenerateTransaction(privA, pubB, 1, nil))
    }

    return pb.(*lazyledger.ProbabilisticBlock), currencyApp.Namespace(), registrarApp.Namespace()
}
