package lazyledger

import (
    "bytes"
    "crypto/sha256"
    "crypto/rand"
    "math"
    "math/big"

    "gitlab.com/NebulousLabs/merkletree"
    "github.com/musalbas/rsmt2d"
)

// ProbabilisticBlock represents a block designed for the Probabilistic Validity Rule.
type ProbabilisticBlock struct {
    prevHash []byte
    messages []Message
    rowRoots [][]byte
    columnRoots [][]byte
    cachedRowRoots [][]byte
    cachedColumnRoots [][]byte
    squareWidth int
    headerOnly bool
    cachedEds *rsmt2d.ExtendedDataSquare
    messageSize int
    validated bool
    sampleRequest *SampleRequest
    provenDependencies map[string]bool
}

type SampleRequest struct {
    Indexes []int
    Axes []int
}

type SampleResponse struct {
    Proofs [][][]byte
}

// NewProbabilisticBlock returns a new probabilistic block.
func NewProbabilisticBlock(prevHash []byte, messageSize int) Block {
    return &ProbabilisticBlock{
        prevHash: prevHash,
        messageSize: messageSize,
        provenDependencies: make(map[string]bool),
    }
}

// ImportProbabilisticBlockBlockHeader imports a received probabilistic block without the messages.
func ImportProbabilisticBlockHeader(prevHash []byte, rowRoots [][]byte, columnRoots [][]byte, squareWidth int, messageSize int, validated bool) Block {
    return &ProbabilisticBlock{
        prevHash: prevHash,
        rowRoots: rowRoots,
        columnRoots: columnRoots,
        squareWidth: squareWidth,
        headerOnly: true,
        messageSize: messageSize,
        validated: validated,
        provenDependencies: make(map[string]bool),
    }
}

// ImportProbabilisticBlock imports a received probabilistic block.
func ImportProbabilisticBlock(prevHash []byte, messages []Message, messageSize int, validated bool) Block {
    return &ProbabilisticBlock{
        prevHash: prevHash,
        messages: messages,
        messageSize: messageSize,
        validated: validated,
        provenDependencies: make(map[string]bool),
    }
}

// SquareWidth returns the width of coded data square of the block.
func (pb *ProbabilisticBlock) SquareWidth() int {
    if pb.headerOnly {
        return pb.squareWidth
    } else {
        return int(pb.eds().Width())
    }
}

// AddMessage adds a message to the block.
func (pb *ProbabilisticBlock) AddMessage(message Message) {
    pb.messages = append(pb.messages, message)
    pb.cachedEds = nil
    pb.cachedRowRoots = nil
    pb.cachedColumnRoots = nil
}

func (pb *ProbabilisticBlock) messagesBytes() [][]byte {
    messagesBytes := make([][]byte, len(pb.messages))
    for index, message := range pb.messages {
        messagesBytes[index] = message.MarshalPadded(pb.messageSize)
    }
    return messagesBytes
}

func (pb *ProbabilisticBlock) eds() *rsmt2d.ExtendedDataSquare {
    if pb.cachedEds == nil {
        data := pb.messagesBytes()
        missingShares := int(math.Pow(math.Ceil(math.Sqrt(float64(len(data)))), 2)) - len(data)
        paddingShare := make([]byte, pb.messageSize)
        for i := 0; i < pb.messageSize; i++ {
            paddingShare[i] = 0xFF // this will ensure it will be treated like a redundancy share
        }
        for i := 0; i < missingShares; i++ {
            freshPaddingShare := make([]byte, pb.messageSize)
            copy(freshPaddingShare, paddingShare)
            data = append(data, freshPaddingShare)
        }
        pb.cachedEds, _ = rsmt2d.ComputeExtendedDataSquare(data, rsmt2d.RSGF8)
    }

    return pb.cachedEds
}

// RowRoots returns the Merkle roots of the rows of the block.
func (pb *ProbabilisticBlock) RowRoots() [][]byte {
    if pb.rowRoots != nil {
        return pb.rowRoots
    }

    if pb.cachedRowRoots == nil {
        pb.computeRoots()
    }

    return pb.cachedRowRoots
}

// ColumnRoots returns the Merkle roots of the columns of the block.
func (pb *ProbabilisticBlock) ColumnRoots() [][]byte {
    if pb.columnRoots != nil {
        return pb.columnRoots
    }

    if pb.cachedColumnRoots == nil {
        pb.computeRoots()
    }

    return pb.cachedColumnRoots
}

func (pb *ProbabilisticBlock) computeRoots() {
    ndf := NewNamespaceDummyFlagger()
    fh := NewFlagHasher(ndf, sha256.New())
    rowRoots := make([][]byte, pb.SquareWidth())
    columnRoots := make([][]byte, pb.SquareWidth())
    var rowTree *merkletree.Tree
    var columnTree *merkletree.Tree
    var rowData [][]byte
    var columnData [][]byte
    for i := 0; i < pb.SquareWidth(); i++ {
        if i >= pb.SquareWidth() / 2 {
            fh.(*flagDigest).setCodedMode(true)
        }
        rowTree = merkletree.New(fh)
        columnTree = merkletree.New(fh)
        rowData = pb.eds().Row(uint(i))
        columnData = pb.eds().Column(uint(i))
        for j := 0; j < pb.SquareWidth(); j++ {
            if j >= pb.SquareWidth() / 2 {
                fh.(*flagDigest).setCodedMode(true)
            }
            rowTree.Push(rowData[j])
            columnTree.Push(columnData[j])
        }
        fh.(*flagDigest).setCodedMode(false)

        rowRoots[i] = rowTree.Root()
        columnRoots[i] = columnTree.Root()
    }

    pb.cachedRowRoots = rowRoots
    pb.cachedColumnRoots = columnRoots
}

func (pb *ProbabilisticBlock) RequestSamples(n int) (*SampleRequest, error) {
    indexes := make([]int, n)
    axes := make([]int, n)
    for i := 0; i < n; i++ {
        val, err := rand.Int(rand.Reader, big.NewInt(int64(math.Pow(float64(pb.SquareWidth()), 2))))
        if err != nil {
            return nil, err
        }
        indexes[i] = int(val.Int64())
        val, err = rand.Int(rand.Reader, big.NewInt(2))
        if err != nil {
            return nil, err
        }
        axes[i] = int(val.Int64())
    }

    pb.sampleRequest = &SampleRequest{
        Indexes: indexes,
        Axes: axes,
    }

    return pb.sampleRequest, nil
}

func (pb *ProbabilisticBlock) RespondSamples(request *SampleRequest) *SampleResponse {
    var proofs [][][]byte
    ndf := NewNamespaceDummyFlagger()
    fh := NewFlagHasher(ndf, sha256.New())
    for x, index := range request.Indexes {
        r, c := pb.shareIndexToCoordinates(index)

        // Add Merkle proof to response
        var data [][]byte
        tree := merkletree.New(fh)
        if request.Axes[x] == 0 { // row
            data = pb.eds().Row(uint(r))
            tree.SetIndex(uint64(c))
            if r >= pb.SquareWidth() / 2 {
                fh.(*flagDigest).setCodedMode(true)
            }
        } else { // column
            data = pb.eds().Column(uint(c))
            tree.SetIndex(uint64(r))
            if c >= pb.SquareWidth() / 2 {
                fh.(*flagDigest).setCodedMode(true)
            }
        }
        for j, share := range data {
            if j >= pb.SquareWidth() / 2 {
                fh.(*flagDigest).setCodedMode(true)
            }
            tree.Push(share)
        }
        fh.(*flagDigest).setCodedMode(false)
        _, proof, _, _ := tree.Prove()
        proofs = append(proofs, proof)
    }

    return &SampleResponse{
        Proofs: proofs,
    }
}

func (pb *ProbabilisticBlock) ProcessSamplesResponse(response *SampleResponse) bool {
    if len(response.Proofs) != len(pb.sampleRequest.Indexes) {
        return false
    }

    ndf := NewNamespaceDummyFlagger()
    fh := NewFlagHasher(ndf, sha256.New())
    for x, index := range pb.sampleRequest.Indexes {
        r, c := pb.shareIndexToCoordinates(index)
        var root []byte
        var result bool
        if r >= pb.SquareWidth() / 2 || c >= pb.SquareWidth() / 2 {
            fh.(*flagDigest).setCodedMode(true)
        }
        if pb.sampleRequest.Axes[x] == 0 { // row
            root = pb.RowRoots()[r]
            result = merkletree.VerifyProof(fh, root, response.Proofs[x], uint64(c), uint64(pb.SquareWidth()))
        } else { // column
            root = pb.ColumnRoots()[c]
            result = merkletree.VerifyProof(fh, root, response.Proofs[x], uint64(r), uint64(pb.SquareWidth()))
        }
        fh.(*flagDigest).setCodedMode(false)
        if !result {
            return false
        }
    }

    return true
}

// Digest computes the hash of the block.
func (pb *ProbabilisticBlock) Digest() []byte {
    hasher := sha256.New()
    hasher.Write(pb.prevHash)
    for _, root := range pb.rowRoots {
        hasher.Write(root)
    }
    for _, root := range pb.columnRoots {
        hasher.Write(root)
    }
    return hasher.Sum(nil)
}

// Valid returns true if the block is valid.
func (pb *ProbabilisticBlock) Valid() bool {
    return pb.validated
}

// PrevHash returns the hash of the previous block.
func (pb *ProbabilisticBlock) PrevHash() []byte {
    return pb.prevHash
}

// Messages returns the block's messages.
func (pb *ProbabilisticBlock) Messages() []Message {
    return pb.messages
}

func (pb *ProbabilisticBlock) indexToCoordinates(index int) (row, column int) {
    row = index / (pb.SquareWidth() / 2)
    column = index % (pb.SquareWidth() / 2)
    return
}

func (pb *ProbabilisticBlock) shareIndexToCoordinates(index int) (row, column int) {
    row = index / pb.SquareWidth()
    column = index % pb.SquareWidth()
    return
}

// ApplicationProof creates a Merkle proof for all of the messages in a block for an application namespace.
// All proofs are created from row roots only.
func (pb *ProbabilisticBlock) ApplicationProof(namespace [namespaceSize]byte) (int, int, [][][]byte, *[]Message, [][]byte) {
    var proofStart int
    var proofEnd int
    var found bool
    for index, message := range pb.messages {
        if message.Namespace() == namespace {
            if !found {
                found = true
                proofStart = index
            }
            proofEnd = index + 1
        }
    }

    var inRange bool
    if !found {
        var prevMessage Message
        // We need to generate a proof for an absence of relevant messages.
        for index, message := range pb.messages {
            if index != 0 {
                prevNs := prevMessage.Namespace()
                currentNs := message.Namespace()
                if ((bytes.Compare(prevNs[:], namespace[:]) < 0 && bytes.Compare(namespace[:], currentNs[:]) < 0) ||
                    (bytes.Compare(prevNs[:], namespace[:]) > 0 && bytes.Compare(namespace[:], currentNs[:]) > 0)) {
                    if !inRange {
                        inRange = true
                        proofStart = index
                    }
                    proofEnd = index + 1
                }
            }
            prevMessage = message
        }
    }

    ndf := NewNamespaceDummyFlagger()
    fh := NewFlagHasher(ndf, sha256.New())
    var proofs [][][]byte
    if found || inRange {
        proofStartRow, proofStartColumn := pb.indexToCoordinates(proofStart)
        proofEndRow, proofEndColumn := pb.indexToCoordinates(proofEnd)
        if proofEndColumn == 0 {
            proofEndRow -= 1
            proofEndColumn = pb.SquareWidth() / 2
        }
        for i := 0; i < pb.SquareWidth() / 2; i++ {
            if i >= proofStartRow && i <= proofEndRow {
                // This row needs Merkle proofs
                var startColumn int
                var endColumn int
                if i == proofStartRow {
                    startColumn = proofStartColumn
                } else {
                    startColumn = 0
                }
                if i == proofEndRow {
                    endColumn = proofEndColumn
                } else {
                    endColumn = pb.SquareWidth() / 2
                }
                rowProof, _ := merkletree.BuildRangeProof(startColumn, endColumn, NewCodedAxisSubtreeHasher(pb.eds().Row(uint(i)), fh))
                proofs = append(proofs, rowProof)
            }
        }
    }
    proofMessages := pb.messages[proofStart:proofEnd]
    if found {
        return proofStart, proofEnd, proofs, &proofMessages, nil
    }

    var hashes [][]byte
    for _, message := range proofMessages {
        ndf := NewNamespaceDummyFlagger()
        fh := NewFlagHasher(ndf, sha256.New())
        hashes = append(hashes, leafSum(fh, message.MarshalPadded(pb.messageSize)))
        fh.Reset()
    }

    return proofStart, proofEnd, proofs, nil, hashes
}

// VerifyApplicationProof verifies a Merkle proof for all of the messages in a block for an application namespace.
func (pb *ProbabilisticBlock) VerifyApplicationProof(namespace [namespaceSize]byte, proofStart int, proofEnd int, proofs [][][]byte, messages *[]Message, hashes [][]byte) bool {
    // Verify Merkle proofs
    ndf := NewNamespaceDummyFlagger()
    fh := NewFlagHasher(ndf, sha256.New())
    var lh merkletree.LeafHasher
    if messages != nil {
        lh = NewPaddedMessageLeafHasher(messages, fh, pb.messageSize)
    } else {
        lh = NewHashLeafHasher(hashes)
    }

    proofStartRow, proofStartColumn := pb.indexToCoordinates(proofStart)
    proofEndRow, proofEndColumn := pb.indexToCoordinates(proofEnd)
    if proofEndColumn == 0 {
        proofEndRow -= 1
        proofEndColumn = pb.SquareWidth() / 2
    }
    proofNum := 0
    for i := 0; i < pb.SquareWidth() / 2; i++ {
        if i >= proofStartRow && i <= proofEndRow {
            // This row has Merkle proofs
            var startColumn int
            var endColumn int
            if i == proofStartRow {
                startColumn = proofStartColumn
            } else {
                startColumn = 0
            }
            if i == proofEndRow {
                endColumn = proofEndColumn
            } else {
                endColumn = pb.SquareWidth() / 2
            }

            // Verify proof
            result, err := merkletree.VerifyRangeProof(lh, fh, startColumn, endColumn, proofs[proofNum], pb.RowRoots()[i])
            if !result || err != nil {
                return false
            }

            // Verify completeness
            var leafIndex uint64
            var leftSubtrees [][]byte
            var rightSubtrees [][]byte
            proof := proofs[proofNum]
        	consumeUntil := func(end uint64) error {
        		for leafIndex != end && len(proof) > 0 {
        			subtreeSize := nextSubtreeSize(leafIndex, end)
                    leftSubtrees = append(leftSubtrees, proof[0])
        			proof = proof[1:]
        			leafIndex += uint64(subtreeSize)
        		}
        		return nil
        	}
            if err := consumeUntil(uint64(startColumn)); err != nil {
                return false
            }
            rightSubtrees = proof

            for _, subtree := range leftSubtrees {
                _, max := dummyNamespacesFromFlag(subtree)
                if bytes.Compare(max, namespace[:]) >= 0 {
                    return false
                }
            }
            for _, subtree := range rightSubtrees {
                min, _ := dummyNamespacesFromFlag(subtree)
                if bytes.Compare(min, namespace[:]) <= 0 {
                    return false
                }
            }

            proofNum += 1
        }
    }

    return true
}

func (pb *ProbabilisticBlock) ProveDependency(index int) ([]byte, [][]byte, error) {
    ndf := NewNamespaceDummyFlagger()
    fh := NewFlagHasher(ndf, sha256.New())
    r, c := pb.indexToCoordinates(index)
    proof, err := merkletree.BuildRangeProof(c, c + 1, NewCodedAxisSubtreeHasher(pb.eds().Row(uint(r)), fh))
    if err != nil {
        return nil, nil, err
    }
    return leafSum(fh, pb.messages[index].MarshalPadded(pb.messageSize)), proof, nil
}

func (pb *ProbabilisticBlock) VerifyDependency(index int, hash []byte, proof [][]byte) bool {
    ndf := NewNamespaceDummyFlagger()
    fh := NewFlagHasher(ndf, sha256.New())
    lh := NewHashLeafHasher([][]byte{hash})
    r, c := pb.indexToCoordinates(index)
    result, err := merkletree.VerifyRangeProof(lh, fh, c, c + 1, proof, pb.RowRoots()[r])
    if result && err == nil {
        pb.provenDependencies[string(hash)] = true
        return true
    }
    return false
}

func (pb *ProbabilisticBlock) DependencyProven(hash []byte) bool {
    if value, ok := pb.provenDependencies[string(hash)]; ok {
        return value
    }
    return false
}
