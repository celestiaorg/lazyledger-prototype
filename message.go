package lazyledger

type Message struct {
    namespace [namespaceSize]byte
    data []byte
}

func NewMessage(namespace [namespaceSize]byte, data []byte) *Message {
    return &Message{
        namespace: namespace,
        data: data,
    }
}

func UnmarshalMessage(marshalled []byte) *Message {
    var namespace [namespaceSize]byte
    copy(namespace[:], marshalled[:namespaceSize])
    return NewMessage(namespace, marshalled[namespaceSize:])
}

func (m *Message) Marshal() []byte {
    return append(m.namespace[:], m.data...)
}
