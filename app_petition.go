package lazyledger

import (
    "encoding/binary"

    "github.com/golang/protobuf/proto"
    "github.com/libp2p/go-libp2p-crypto"
)

type PetitionApp struct {
    state MapStore
}

func NewPetitionApp(state MapStore) Application {
    return &PetitionApp{
        state: state,
    }
}

func (app *PetitionApp) ProcessMessage(message Message) {
    transaction := &PetitionAppTransaction{}
    err := proto.Unmarshal(message.Data(), transaction)
    if err != nil {
        return
    }
    apm := transaction.GetApm()
    if apm != nil {
        app.ProcessAddPetitionMessage(apm)
    }
    spm := transaction.GetSpm()
    if spm != nil {
        app.ProcessSignPetitionMessage(spm)
    }
}

func (app *PetitionApp) ProcessAddPetitionMessage(apm *AddPetitionMessage) {
    app.addPetition(*apm.Text)
}

func (app *PetitionApp) ProcessSignPetitionMessage(spm *SignPetitionMessage) {
    key, _ := crypto.UnmarshalPublicKey(spm.Signer)
    signedData := make([]byte, binary.MaxVarintLen64)
    binary.BigEndian.PutUint64(signedData, *spm.Id)
    ok, err := key.Verify(signedData, spm.Signature)
    if !ok || err != nil {
        return
    }
    app.incrementPetition(*spm.Id)
}

func (app *PetitionApp) Namespace() [namespaceSize]byte {
    var namespace [namespaceSize]byte
    copy(namespace[:], []byte("pet"))
    return namespace
}

func (app *PetitionApp) SetBlockHead(hash []byte) {
    app.state.Put([]byte("__head__"), hash)
}

func (app *PetitionApp) BlockHead() []byte {
    head, _ := app.state.Get([]byte("__head__"))
    return head
}

func (app *PetitionApp) Petition(petition uint64) uint64 {
    id := make([]byte, binary.MaxVarintLen64)
    binary.BigEndian.PutUint64(id, petition)
    c, err := app.state.Get(id)
    if err != nil {
        return 0
    }
    return binary.BigEndian.Uint64(c)
}

func (app *PetitionApp) incrementPetition(petition uint64) {
    v := app.Petition(petition)
    newValue := make([]byte, binary.MaxVarintLen64)
    binary.BigEndian.PutUint64(newValue, v + 1)

    id := make([]byte, binary.MaxVarintLen64)
    binary.BigEndian.PutUint64(id, petition)
    app.state.Put(id, newValue)
}

func (app *PetitionApp) addPetition(text string) {
    latestIdBytes, err := app.state.Get([]byte("__last__"))
    var latestId uint64
    if err != nil {
        latestId = 0
    } else {
        latestId = binary.BigEndian.Uint64(latestIdBytes)
    }
    newValue := make([]byte, binary.MaxVarintLen64)
    binary.BigEndian.PutUint64(newValue, latestId + 1)
    app.state.Put([]byte("__last__"), newValue)
    app.state.Put(append([]byte("text__"), newValue...), []byte(text))
}

func (app *PetitionApp) GenerateSignPetitionTransaction(key crypto.PrivKey, petition uint64) Message {
    signedData := make([]byte, binary.MaxVarintLen64)
    binary.BigEndian.PutUint64(signedData, petition)
    sig, _ := key.Sign(signedData)
    kb, _ := key.GetPublic().Bytes()
    spm := &SignPetitionMessage{
        Id: &petition,
        Signature: sig,
        Signer: kb,
    }
    pspm := PetitionAppTransaction_Spm{Spm: spm}
    t := &PetitionAppTransaction{
        Message: &pspm,
    }
    d, _ := proto.Marshal(t)
    return *NewMessage(app.Namespace(), d)
}

func (app *PetitionApp) GenerateAddPetitionTransaction(text string) Message {
    apm := &AddPetitionMessage{
        Text: &text,
    }
    papm := PetitionAppTransaction_Apm{Apm: apm}
    t := &PetitionAppTransaction{
        Message: &papm,
    }
    d, _ := proto.Marshal(t)
    return *NewMessage(app.Namespace(), d)
}
