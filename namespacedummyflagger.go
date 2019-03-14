package lazyledger

import (
    "bytes"
)

type namespaceDummyFlagger struct {}

// NewNamespaceDummyFlagger returns a new dummy flagger for namespaced Merkle trees.
func NewNamespaceDummyFlagger() Flagger {
    return &namespaceDummyFlagger{}
}

func (namespaceDummyFlagger) LeafFlag(leaf []byte) []byte {
    return append(leaf[:namespaceSize], leaf[:namespaceSize]...)
}

func (namespaceDummyFlagger) NodeFlag(node []byte) []byte {
    return node[:flagSize]
}

func (namespaceDummyFlagger) Union(leftFlag []byte, rightFlag []byte) []byte {
    namespaces := make([][]byte, 4)
    namespaces[0], namespaces[1] = dummyNamespacesFromFlag(leftFlag)
    namespaces[2], namespaces[3] = dummyNamespacesFromFlag(rightFlag)

    minNamespace := namespaces[0]
    maxNamespace := namespaces[0]
    for _, namespace := range namespaces[1:] {
        if bytes.Compare(minNamespace, namespace) > 0 {
            minNamespace = namespace
        }
        if bytes.Compare(maxNamespace, namespace) < 0 {
            maxNamespace = namespace
        }
    }

    return dummyFlagFromNamespaces(minNamespace, maxNamespace)
}

func (namespaceDummyFlagger) FlagSize() int {
    return flagSize
}

func dummyNamespacesFromFlag(flag []byte) ([]byte, []byte) {
    return flag[:namespaceSize], flag[namespaceSize:flagSize]
}

func dummyFlagFromNamespaces(leftNamespace []byte, rightNamespace []byte) []byte {
    return append(leftNamespace, rightNamespace...)
}
