package lazyledger

import (
    "crypto/rand"
    "testing"

    "github.com/libp2p/go-libp2p-crypto"
)

func TestAppPetitionSimpleBlock(t *testing.T) {
    bs := NewSimpleBlockStore()
    b := NewBlockchain(bs)

    sb := NewSimpleBlock([]byte{0})

    ms := NewSimpleMap()
    app := NewPetitionApp(ms)
    b.RegisterApplication(&app)

    privA, _, _ := crypto.GenerateSecp256k1Key(rand.Reader)

    sb.AddMessage(app.(*PetitionApp).GenerateAddPetitionTransaction("foo"))
    sb.AddMessage(app.(*PetitionApp).GenerateSignPetitionTransaction(privA, 0))
    b.ProcessBlock(sb)

    if app.(*PetitionApp).Petition(0) != 1 {
        t.Error("failed to sign petition")
    }
}
