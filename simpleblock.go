package lazyledger

// SimpleBlock represents a block designed for the Simple Validity Rule.
type SimpleBlock struct {
    prevHash []byte
    messages []Message
}

// NewSimpleBlock returns a new simple block.
func NewSimpleBlock(prevHash []byte) Block {
    return &SimpleBlock{
        prevHash: prevHash,
    }
}

// AddMessage adds a message to the block.
func (sb *SimpleBlock) AddMessage(message Message) {
    sb.messages = append(sb.messages, message)
}
