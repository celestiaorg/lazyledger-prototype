package main

import (
    "crypto/rand"
    "encoding/binary"
    "fmt"

    "github.com/lazyledger/lazyledger-prototype"
    "github.com/libp2p/go-libp2p-crypto"
)

var privA crypto.PrivKey
var pubA crypto.PubKey
var pubB crypto.PubKey
var pubC crypto.PubKey

const namespaceSize = 8

func main() {
    privA, pubA, _ = crypto.GenerateSecp256k1Key(rand.Reader)
    _, pubB, _ = crypto.GenerateSecp256k1Key(rand.Reader)
    _, pubC, _ = crypto.GenerateSecp256k1Key(rand.Reader)

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
            _, _, proofs1, messages1, _ = sb.ApplicationProof(rns)
            for _, msg := range *messages1 {
                sbBandwidth += len(msg.Marshal())
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
            _, _, proofs2, messages2, _ = pb.ApplicationProof(rns)
            for _, msg := range *messages2 {
                pbBandwidth += len(msg.Marshal())
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

    pubABytes, _ := pubA.Bytes()
    pubBBytes, _ := pubB.Bytes()
    pubCBytes, _ := pubC.Bytes()
    pubABalanceBytes := make([]byte, binary.MaxVarintLen64)
    binary.BigEndian.PutUint64(pubABalanceBytes, 1000000)
    ms1.Put(pubABytes, pubABalanceBytes)

    ms2 := lazyledger.NewSimpleMap()
    registrarApp := lazyledger.NewRegistrar(ms2, currencyApp.(*lazyledger.Currency), pubBBytes)
    var rns1 [namespaceSize]byte
    copy(rns1[:], []byte("reg1"))
    registrarApp.(*lazyledger.Registrar).SetNamespace(rns1)
    b.RegisterApplication(&registrarApp)
    sb.AddMessage(currencyApp.(*lazyledger.Currency).GenerateTransaction(privA, pubB, 100000, nil))

    for i := 0; i < registrarTxes; i++ {
        name := make([]byte, 8)
        //rand.Read(name)
        sb.AddMessage(registrarApp.(*lazyledger.Registrar).GenerateTransaction(privA, name))
    }

    ms3 := lazyledger.NewSimpleMap()
    registrarApp2 := lazyledger.NewRegistrar(ms3, currencyApp.(*lazyledger.Currency), pubCBytes)
    var rns2 [namespaceSize]byte
    copy(rns2[:], []byte("reg2"))
    registrarApp2.(*lazyledger.Registrar).SetNamespace(rns2)
    b.RegisterApplication(&registrarApp2)
    sb.AddMessage(currencyApp.(*lazyledger.Currency).GenerateTransaction(privA, pubC, 100000, nil))

    for i := 0; i < otherTxes; i++ {
        name := make([]byte, 8)
        //rand.Read(name)
        sb.AddMessage(registrarApp2.(*lazyledger.Registrar).GenerateTransaction(privA, name))
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

    pubABytes, _ := pubA.Bytes()
    pubBBytes, _ := pubB.Bytes()
    pubCBytes, _ := pubC.Bytes()
    pubABalanceBytes := make([]byte, binary.MaxVarintLen64)
    binary.BigEndian.PutUint64(pubABalanceBytes, 1000000)
    ms1.Put(pubABytes, pubABalanceBytes)

    ms2 := lazyledger.NewSimpleMap()
    registrarApp := lazyledger.NewRegistrar(ms2, currencyApp.(*lazyledger.Currency), pubBBytes)
    var rns1 [namespaceSize]byte
    copy(rns1[:], []byte("reg1"))
    registrarApp.(*lazyledger.Registrar).SetNamespace(rns1)
    b.RegisterApplication(&registrarApp)
    pb.AddMessage(currencyApp.(*lazyledger.Currency).GenerateTransaction(privA, pubB, 100000, nil))

    for i := 0; i < registrarTxes; i++ {
        name := make([]byte, 8)
        //rand.Read(name)
        pb.AddMessage(registrarApp.(*lazyledger.Registrar).GenerateTransaction(privA, name))
    }

    ms3 := lazyledger.NewSimpleMap()
    registrarApp2 := lazyledger.NewRegistrar(ms3, currencyApp.(*lazyledger.Currency), pubCBytes)
    var rns2 [namespaceSize]byte
    copy(rns2[:], []byte("reg2"))
    registrarApp2.(*lazyledger.Registrar).SetNamespace(rns2)
    b.RegisterApplication(&registrarApp2)
    pb.AddMessage(currencyApp.(*lazyledger.Currency).GenerateTransaction(privA, pubC, 100000, nil))

    for i := 0; i < otherTxes; i++ {
        name := make([]byte, 8)
        //rand.Read(name)
        pb.AddMessage(registrarApp2.(*lazyledger.Registrar).GenerateTransaction(privA, name))
    }

    return pb.(*lazyledger.ProbabilisticBlock), currencyApp.Namespace(), registrarApp.Namespace()
}
