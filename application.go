package lazyledger

// Application is an interface for a lazyledger application.
type Application interface {
    // ProcessMessage processes a message according to the application's state machine.
    ProcessMessage(message Message) bool

    // Namespace returns the namespace ID of the application.
    Namespace() [namespaceSize]byte

    // BlockHead returns the hash of the latest block that has been processed.
    BlockHead() []byte

    // SetBlockHead sets the hash of the latest block that has been processed.
    SetBlockHead(hash []byte)
}
