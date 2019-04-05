package lazyledger

const namespaceSize = 8
const flagSize = 16

var codedNamespace [namespaceSize]byte
var codedFlag [flagSize]byte

func init() {
    for i, _ := range codedNamespace {
        codedNamespace[i] = 0xFF
    }
    for i, _ := range codedFlag {
        codedFlag[i] = 0xFF
    }
}
