package merkle

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"log"
	"math"
	"os"
)

type MerkleNode struct {
	parent     *MerkleNode
	LeftChild  *MerkleNode
	RightChild *MerkleNode
	IsLeaf     bool
	Val        string

	NodeHash []byte
	Tag      string
}

type MerkleTree struct {
	RootNode  *MerkleNode
	RootHash  []byte
	N         int
	LeafNodes []MerkleNode
}

func (n *MerkleNode) generateNodeHash() {

	if n.IsLeaf {
		out := sha256.Sum256([]byte(n.Val))
		n.NodeHash = out[:]
	} else {
		if n.RightChild == nil {
			n.NodeHash = n.LeftChild.NodeHash
		} else {
			// Refer: https://stackoverflow.com/questions/59860517/calculating-sha256-gives-different-results-after-appending-slices-depending-on-i
			var hashSlice []byte
			hashSlice = append(hashSlice, n.LeftChild.NodeHash...)
			hashSlice = append(hashSlice, n.RightChild.NodeHash...)

			// parent hash =  hash (child1 | child2)
			temp := sha256.Sum256(hashSlice)
			n.NodeHash = temp[:]
		}
	}
}

func (n *MerkleNode) BuildTree(arr []string) {

	if len(arr) == 1 {
		n.IsLeaf = true
		n.LeftChild = nil
		n.RightChild = nil
		n.parent = n
		n.Val = arr[0]

		n.generateNodeHash()
		return
	}

	pow2 := int(math.Floor(math.Log2(float64(len(arr)))))
	mid := 1 << pow2

	if mid == len(arr) {
		pow2--
		mid = 1 << pow2
	}

	n.LeftChild = NewNode()
	n.LeftChild.parent = &MerkleNode{}
	n.LeftChild.BuildTree(arr[:mid])

	n.RightChild = NewNode()
	n.RightChild.parent = n
	n.RightChild.BuildTree(arr[mid:])

	n.generateNodeHash()
}

func (n *MerkleNode) MarshalJSON() ([]byte, error) {
	intermediate := map[string]interface{}{
		"LeftChild":  n.LeftChild,
		"RightChild": n.RightChild,
		"IsLeaf":     n.IsLeaf,
		"Val":        n.Val,
		"NodeHash":   hex.EncodeToString(n.NodeHash),
		"Tag":        n.Tag,
	}
	return json.Marshal(intermediate)
}

func (m *MerkleTree) SaveToJSON() {
	bytes, err := json.MarshalIndent(m, "", "\t")
	if err != nil {
		log.Fatal("Error cannot convert tree to JSON", err)
	}

	f, err := os.Create(hex.EncodeToString(m.RootHash))
	if err != nil {
		log.Fatal("Error cannot create output file")
	}
	defer f.Close()

	if _, err := f.Write(bytes); err != nil {
		panic(err)
	}
}

func NewNode() *MerkleNode {
	return &MerkleNode{
		IsLeaf:     false,
		LeftChild:  nil,
		RightChild: nil,
		parent:     nil,
		Val:        "",
		NodeHash:   []byte{},
		Tag:        "",
	}
}

func NewTree(arr []string) *MerkleTree {
	tree := &MerkleTree{
		RootNode:  NewNode(),
		RootHash:  []byte{},
		N:         len(arr),
		LeafNodes: []MerkleNode{},
	}

	tree.RootNode.BuildTree(arr)
	tree.RootHash = tree.RootNode.NodeHash

	return tree
}

// To-Do
func (m *MerkleTree) verifyRootHash() {

}
