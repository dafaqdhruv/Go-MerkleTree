package merkle

import (
	"crypto/sha256"
	"math"
)

type MerkleNode struct {
	parent     *MerkleNode
	leftChild  *MerkleNode
	rightChild *MerkleNode
	isLeaf     bool
	val        string

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

	if n.isLeaf {
		out := sha256.Sum256([]byte(n.val))
		n.NodeHash = out[:]
	} else {
		if n.rightChild == nil {
			n.NodeHash = n.leftChild.NodeHash
		} else {
			// Refer: https://stackoverflow.com/questions/59860517/calculating-sha256-gives-different-results-after-appending-slices-depending-on-i
			var hashSlice []byte
			hashSlice = append(hashSlice, n.leftChild.NodeHash...)
			hashSlice = append(hashSlice, n.rightChild.NodeHash...)

			// parent hash =  hash (child1 | child2)
			temp := sha256.Sum256(hashSlice)
			n.NodeHash = temp[:]
		}
	}
}

func (n *MerkleNode) BuildTree(arr []string) {

	if len(arr) == 1 {
		n.isLeaf = true
		n.leftChild = nil
		n.rightChild = nil
		n.parent = n
		n.val = arr[0]

		n.generateNodeHash()
		return
	}

	pow2 := int(math.Floor(math.Log2(float64(len(arr)))))
	mid := 1 << pow2

	if mid == len(arr) {
		pow2--
		mid = 1 << pow2
	}

	n.leftChild = NewNode()
	n.leftChild.parent = &MerkleNode{}
	n.leftChild.BuildTree(arr[:mid])

	n.rightChild = NewNode()
	n.rightChild.parent = n
	n.rightChild.BuildTree(arr[mid:])

	n.generateNodeHash()
}

func NewNode() *MerkleNode {
	return &MerkleNode{
		isLeaf:     false,
		leftChild:  nil,
		rightChild: nil,
		parent:     nil,
		val:        "",
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
