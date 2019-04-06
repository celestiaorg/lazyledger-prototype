package lazyledger

import (
    "github.com/golang/protobuf/proto"
)

type DummyApp struct {
    state MapStore
}

func NewDummyApp(state MapStore) Application {
    return &DummyApp{
        state: state,
    }
}

func (app *DummyApp) ProcessMessage(message Message) {
    transaction := &DummyAppTransaction{}
    err := proto.Unmarshal(message.Data(), transaction)
    if err != nil {
        return
    }
    for k, v := range transaction.Puts {
        app.state.Put([]byte(k), []byte(v))
    }
}

func (app *DummyApp) Namespace() [namespaceSize]byte {
    var namespace [namespaceSize]byte
    copy(namespace[:], []byte("dummy"))
    return namespace
}

func (app *DummyApp) SetBlockHead(hash []byte) {
    app.state.Put([]byte("__head__"), hash)
}

func (app *DummyApp) BlockHead() []byte {
    head, _ := app.state.Get([]byte("__head__"))
    return head
}

func (app *DummyApp) Get(key string) string {
    value, err := app.state.Get([]byte(key))
    if err != nil {
        return ""
    }
    return string(value)
}

func (app *DummyApp) GenerateTransaction(puts map[string]string) Message {
    transaction := &DummyAppTransaction{
        Puts: puts,
    }
    data, _ := proto.Marshal(transaction)
    return *NewMessage(app.Namespace(), data)
}
