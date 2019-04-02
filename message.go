package lazyledger

import (
    "encoding/binary"
)

// Message represents a namespaced message.
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

// UnmarshalPaddedMessage returns a message from its marshalled padded raw data.
func UnmarshalPaddedMessage(marshalled []byte) *Message {
    marshalledSizeBytes := make([]byte, 2)
    marshalledSizeBytes[0] = marshalled[len(marshalled) - 2]
    marshalledSizeBytes[1] = marshalled[len(marshalled) - 1]
    marshalled = marshalled[:int(binary.LittleEndian.Uint16(marshalledSizeBytes))]

    var namespace [namespaceSize]byte
    copy(namespace[:], marshalled[:namespaceSize])
    return NewMessage(namespace, marshalled[namespaceSize:])
}

// Marshal converts a message to raw data.
func (m *Message) Marshal() []byte {
    return append(m.namespace[:], m.data...)
}

// Marshal converts a message to padded raw data.
func (m *Message) MarshalPadded(messageSize int) []byte {
    marshalled := append(m.namespace[:], m.data...)
    marshalledSizeBytes := make([]byte, 2)
    binary.LittleEndian.PutUint16(marshalledSizeBytes, uint16(len(marshalled)))

    padding := make([]byte, messageSize - len(marshalled))
    for i, _ := range padding {
        padding[i] = 0x00
    }
    marshalled = append(marshalled, padding...)
    marshalled[len(marshalled) - 2] = marshalledSizeBytes[0]
    marshalled[len(marshalled) - 1] = marshalledSizeBytes[1]
    return marshalled
}

// Namespace returns the namespace of a message.
func (m *Message) Namespace() [namespaceSize]byte {
    return m.namespace;
}

// Data returns the data of a message.
func (m *Message) Data() []byte {
    return m.data;
}
