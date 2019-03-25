package lazyledger

import (
    "testing"
)

func TestSimpleBlock(t *testing.T) {
    sb := NewSimpleBlock([]byte{0})

    sb.AddMessage(*NewMessage([namespaceSize]byte{0}, []byte("foo")))
    sb.AddMessage(*NewMessage([namespaceSize]byte{1}, []byte("foo")))
    sb.AddMessage(*NewMessage([namespaceSize]byte{1}, []byte("foo")))
    sb.AddMessage(*NewMessage([namespaceSize]byte{2}, []byte("foo")))

    proofStart, proofEnd, proof, messages := sb.ApplicationProof([namespaceSize]byte{1})
    result := sb.VerifyApplicationProof([namespaceSize]byte{1}, proofStart, proofEnd, proof, messages)
    if !result {
        t.Error("VerifyApplicationProof incorrectly returned false")
    }

    // TODO: add negative tests
}
