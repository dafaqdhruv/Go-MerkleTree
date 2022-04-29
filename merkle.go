package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
)

type node struct {
	nodeHash []byte
	parent   *node
	left     *node
	right    *node
	isLeaf   bool
	val      string
}

type merkleTree struct {
	root      *node
	rootHash  []byte
	n         int
	leafNodes []node
}

func (n *node) getNodeHash() {

	// if Leaf node : get hash of value
	if n.isLeaf {
		out := sha256.Sum256([]byte(n.val))
		n.nodeHash = out[:]
	} else {

		// if node contains only 1 child
		// propagate the same hash upwards
		if n.right == nil && n.left != nil {
			n.nodeHash = n.left.nodeHash
		} else {

			// else node contains 2 children
			// parent hash =  hash (child1 | child2)

			// ISSUE : sha256 gives different results of hashing.
			// Solved : https://stackoverflow.com/questions/59860517/calculating-sha256-gives-different-results-after-appending-slices-depending-on-i

			var hashSlice []byte
			hashSlice = append(hashSlice, n.left.nodeHash...)
			hashSlice = append(hashSlice, n.right.nodeHash...)

			temp := sha256.Sum256(hashSlice)
			n.nodeHash = temp[:]
		}
	}
}

// takes data array input and returns a merkle tree
func buildTree(arr []string) *merkleTree {

	n := len(arr)

	outTree := new(merkleTree)
	outTree.n = n
	outTree.leafNodes = make([]node, n)

	// fill in leaf nodes
	for i, j := range arr {

		outTree.leafNodes[i] = node{
			left:   nil,
			right:  nil,
			isLeaf: true,
			val:    j,
		}
		outTree.leafNodes[i].getNodeHash()
	}

	tempSlice := outTree.leafNodes // slice to store nodes temporarily

	lim := n // number of elements in the level below
	height := math.Ceil(math.Log2(float64(n)))

	// iterate level wise
	for level := height - 1; level >= 0; level-- {

		tempSlice2 := make([]node, n)
		elems := 0 // will count elements in this level

		for nodex := 0; nodex < lim; nodex += 2 {

			// fmt.Printf("working on %f level and %d node of limit %d\n", level, nodex, lim)
			// if there are odd number of elements in a level
			if lim%2 == 1 && nodex == lim-1 {
				tempParent := new(node)
				tempParent.left = &tempSlice[nodex]
				tempParent.right = nil
				tempParent.isLeaf = false

				tempSlice[nodex].parent = tempParent
				tempParent.getNodeHash()
				tempSlice2[nodex/2] = *tempParent
				elems++

				if level == 0 {
					outTree.root = tempParent
					outTree.rootHash = tempParent.nodeHash

				}
			} else {

				nodey := nodex + 1 // select pairwise elements

				tempParent := new(node)
				tempParent.left = &tempSlice[nodex]
				tempParent.right = &tempSlice[nodey]
				tempParent.isLeaf = false

				tempSlice[nodex].parent = tempParent
				tempSlice[nodey].parent = tempParent

				tempParent.getNodeHash()
				tempSlice2[nodex/2] = *tempParent
				elems += 1

				// debug msg
				// fmt.Println(hex.EncodeToString(tempParent.nodeHash))

				if level == 0 {
					outTree.root = tempParent
					outTree.rootHash = tempParent.nodeHash
				}
			}
		}
		tempSlice = tempSlice2
		lim = elems

	}
	return outTree
}

// To-Do
func (m *merkleTree) verifyRootHash() {

}

func main() {

	words := []string{"hello", "this", "is", "a", "merkle", "tree"}

	Mtree := buildTree(words)

	fmt.Println("root hash is ")
	fmt.Println(hex.EncodeToString(Mtree.rootHash))

	// debug
	// nodeHashes of left and right nodes
	fmt.Println("left hash is ")
	fmt.Println(hex.EncodeToString(Mtree.root.left.nodeHash))

	fmt.Println("right hash is")
	fmt.Println(hex.EncodeToString(Mtree.root.right.nodeHash))

}
