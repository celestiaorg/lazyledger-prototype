package lazyledger

type Message struct {
    namespace [32]byte
    data []byte
}

func NewMessage(namespace [32]byte, data []byte) *Message {
    return &Message{
        namespace: namespace,
        data: data,
    }
}
