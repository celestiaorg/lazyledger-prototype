package main

import (
    "fmt"
    "crypto/rand"

    "github.com/musalbas/lazyledger-prototype"
)

const namespaceSize = 8

func main() {
    txAmounts := []int{128, 256, 384, 512, 640, 1024}
    txSize := 128
    for _, txes := range txAmounts {
        sb := generateSimpleBlock(txes, txSize)
        sbBandwidth := len(sb.PrevHash()) + len(sb.MessagesRoot()) + txSize * txes

        pb := generateProbabilisticBlock(txes, txSize)
        req, _ := pb.RequestSamples(10)
        res := pb.RespondSamples(req)
        pbBandwidth := 0
        for _, root := range pb.RowRoots() {
            pbBandwidth += len(root)
        }
        for _, root := range pb.ColumnRoots() {
            pbBandwidth += len(root)
        }
        for _, proof := range res.Proofs {
            for _, hash := range proof {
                pbBandwidth += len(hash)
            }
        }

        fmt.Println(txes, sbBandwidth, pbBandwidth)
    }
}

func generateSimpleBlock(txes int, txSize int) *lazyledger.SimpleBlock {
    txSize -= namespaceSize
    sb := lazyledger.NewSimpleBlock([]byte{0})

    for i := 0; i < txes; i++ {
        messageData := make([]byte, txSize)
        rand.Read(messageData)
        sb.AddMessage(*lazyledger.NewMessage([namespaceSize]byte{0}, messageData))
    }

    return sb.(*lazyledger.SimpleBlock)
}

func generateProbabilisticBlock(txes int, txSize int) *lazyledger.ProbabilisticBlock {
    pb := lazyledger.NewProbabilisticBlock([]byte{0}, txSize)
    txSize -= namespaceSize + 2

    for i := 0; i < txes; i++ {
        messageData := make([]byte, txSize)
        rand.Read(messageData)
        pb.AddMessage(*lazyledger.NewMessage([namespaceSize]byte{0}, messageData))
    }

    return pb.(*lazyledger.ProbabilisticBlock)
}
