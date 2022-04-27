package main

import (
	"crypto/sha256"
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

	if n.isLeaf {

		out := sha256.Sum256([]byte(n.val))
		n.nodeHash = out[:]
	} else {

		// if node contains only 1 child
		// transfer the same hash upwards

		if n.right == nil {

			n.nodeHash = n.left.nodeHash
		}

		// else node contains 2 children
		// parent hash =  hash (child1 | child2)
		hashSlice := []byte{}
		hashSlice = append(hashSlice, n.left.nodeHash...)
		hashSlice = append(hashSlice, n.right.nodeHash...)

		temp := sha256.Sum256(hashSlice)
		n.nodeHash = temp[:]

	}
}

func buildTree(arr []string) *merkleTree {

	n := len(arr)

	// height of tree = ceil (log(n))
	// max nodes in tree of height h =

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

	//
	// iterate level wise
	tempSlice := outTree.leafNodes

	for level := math.Ceil(math.Log2(float64(n))) - 1; level > 0; level-- {

		for nodex := 0; nodex <= int(math.Pow(2, level-1)); nodex += 2 {

			nodey := nodex + 1

			tempParent := &node{
				left:   &tempSlice[nodex],
				right:  &tempSlice[nodey],
				isLeaf: false,
			}

			tempSlice[nodex].parent = tempParent
			tempSlice[nodey].parent = tempParent

			tempParent.getNodeHash()
			tempSlice[nodex/2] = *tempParent

			if level == 1 {
				outTree.root = tempParent
				outTree.rootHash = tempParent.nodeHash
			}
		}
	}
	fmt.Println(outTree.root)
	fmt.Println((outTree.rootHash))
	return outTree
}

func (m *merkleTree) verifyRootHash() {

}

func main() {

	words := []string{"hello", "this", "is", "a", "merkle", "tree"}

	bruh := buildTree(words)

	fmt.Println((bruh.rootHash))
	fmt.Println(bruh.root)
	// fmt.Println(bruh.root.right)
	// fmt.Println(bruh.root.left)

	tempSlices := []byte{}
	for _, x := range words {
		tempSlices = append(tempSlices, []byte(x)...)
	}

	tempHash := sha256.Sum256(tempSlices)
	fmt.Println(tempHash)

}
