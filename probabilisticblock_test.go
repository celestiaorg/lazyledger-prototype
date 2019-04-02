package lazyledger

import (
    "testing"
)

func TestProbabilisticBlock(t *testing.T) {
    pb := NewProbabilisticBlock([]byte{0}, 512)

    pb.AddMessage(*NewMessage([namespaceSize]byte{0}, []byte("foo")))
    pb.AddMessage(*NewMessage([namespaceSize]byte{1}, []byte("foo")))
    pb.AddMessage(*NewMessage([namespaceSize]byte{1}, []byte("foo")))
    pb.AddMessage(*NewMessage([namespaceSize]byte{1}, []byte("foo")))
    pb.AddMessage(*NewMessage([namespaceSize]byte{3}, []byte("foo")))
    pb.AddMessage(*NewMessage([namespaceSize]byte{3}, []byte("foo")))
    pb.AddMessage(*NewMessage([namespaceSize]byte{4}, []byte("foo")))
    pb.AddMessage(*NewMessage([namespaceSize]byte{4}, []byte("foob")))

    proofStart, proofEnd, proofs, messages, hashes := pb.(*ProbabilisticBlock).ApplicationProof([namespaceSize]byte{1})
    if messages == nil {
        t.Error("ApplicationProof incorrectly returned no messages")
    }
    result := pb.(*ProbabilisticBlock).VerifyApplicationProof([namespaceSize]byte{1}, proofStart, proofEnd, proofs, messages, hashes)
    if !result {
        t.Error("VerifyApplicationProof incorrectly returned false")
    }

    proofStart, proofEnd, proofs, messages, hashes = pb.(*ProbabilisticBlock).ApplicationProof([namespaceSize]byte{1})
    proofs[0][0][0] = 0xFF
    if messages == nil {
        t.Error("ApplicationProof incorrectly returned no messages")
    }
    result = pb.(*ProbabilisticBlock).VerifyApplicationProof([namespaceSize]byte{1}, proofStart, proofEnd, proofs, messages, hashes)
    if result {
        t.Error("VerifyApplicationProof incorrectly returned true")
    }

    proofStart, proofEnd, proofs, messages, hashes = pb.(*ProbabilisticBlock).ApplicationProof([namespaceSize]byte{2})
    if messages != nil {
        t.Error("ApplicationProof incorrectly returned messages")
    }
    result = pb.(*ProbabilisticBlock).VerifyApplicationProof([namespaceSize]byte{2}, proofStart, proofEnd, proofs, messages, hashes)
    if !result {
        t.Error("VerifyApplicationProof incorrectly returned false")
    }

    proofStart, proofEnd, proofs, messages, hashes = pb.(*ProbabilisticBlock).ApplicationProof([namespaceSize]byte{2})
    proofs[0][0][0] = 0xFF
    if messages != nil {
        t.Error("ApplicationProof incorrectly returned messages")
    }
    result = pb.(*ProbabilisticBlock).VerifyApplicationProof([namespaceSize]byte{2}, proofStart, proofEnd, proofs, messages, hashes)
    if result {
        t.Error("VerifyApplicationProof incorrectly returned true")
    }

    // TODO: add negative tests
}
