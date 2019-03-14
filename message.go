package lazyledger

type Message struct {
    namespace [namespaceSize]byte
    data []byte
}

// NewMessage returns a new message from its namespace and data.
func NewMessage(namespace [namespaceSize]byte, data []byte) *Message {
    return &Message{
        namespace: namespace,
        data: data,
    }
}

// UnmarshalMessage returns a message from its marshalled raw data.
func UnmarshalMessage(marshalled []byte) *Message {
    var namespace [namespaceSize]byte
    copy(namespace[:], marshalled[:namespaceSize])
    return NewMessage(namespace, marshalled[namespaceSize:])
}

// Marshal converts a message to raw data.
func (m *Message) Marshal() []byte {
    return append(m.namespace[:], m.data...)
}

// Namespace returns the namespace of a message.
func (m *Message) Namespace() [namespaceSize]byte {
    return m.namespace;
}

// Data returns the data of a message.
func (m *Message) Data() []byte {
    return m.data;
}
