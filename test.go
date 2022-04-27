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
			fmt.Println("hashCOPY")
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
	// fmt.Println(hex.EncodeToString(n.nodeHash) + " for " + n.val)

	n := len(arr)
	// height of tree = ceil (log(n))

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

	// iterate level wise
	tempSlice := outTree.leafNodes

	lim := n
	height := math.Ceil(math.Log2(float64(n)))
	for level := height - 1; level >= 0; level-- {

		elems := 0
		for nodex := 0; nodex < lim; nodex += 2 {

			if lim < 2 {
				break
			}
			if lim-nodex < 2 {
				tempParent := &node{
					left:   &tempSlice[nodex],
					right:  nil,
					isLeaf: false,
				}
				tempSlice[nodex].parent = tempParent
				tempParent.getNodeHash()
				tempSlice[nodex/2] = *tempParent
				elems++

				if level <= 0 {
					outTree.root = tempParent
					outTree.rootHash = tempParent.nodeHash

				}
			} else {

				nodey := nodex + 1 // select pairwise elements

				tempParent := &node{
					left:   &tempSlice[nodex],
					right:  &tempSlice[nodey],
					isLeaf: false,
				}

				tempSlice[nodex].parent = tempParent
				tempSlice[nodey].parent = tempParent

				tempParent.getNodeHash()
				tempSlice[nodex/2] = *tempParent
				elems += 1

				if level <= 0 {
					outTree.root = tempParent
					outTree.rootHash = tempParent.nodeHash
				}
			}
		}
		fmt.Print("nodes this level are : ")
		fmt.Println(lim)
		lim = elems

	}

	nodex := 0
	nodey := 1 // select pairwise elements

	tempParent := &node{
		left:   &tempSlice[nodex],
		right:  &tempSlice[nodey],
		isLeaf: false,
	}

	tempSlice[nodex].parent = tempParent
	tempSlice[nodey].parent = tempParent

	tempParent.getNodeHash()
	tempSlice[nodex/2] = *tempParent

	outTree.root = tempParent
	outTree.rootHash = tempParent.nodeHash

	// fmt.Println((outTree.rootHash))
	return outTree
}

func (m *merkleTree) verifyRootHash() {

}

func main() {

	words := []string{"hello", "this", "is", "a", "merkle", "tree"}

	Mtree := buildTree(words)

	fmt.Println(hex.EncodeToString(Mtree.rootHash))

	// nodeHashes of left and right nodes
	fmt.Println("root hash is ")
	fmt.Println(hex.EncodeToString(Mtree.rootHash))

	fmt.Println("left hash is ")
	fmt.Println(hex.EncodeToString(Mtree.root.left.nodeHash))

	fmt.Println("right hash is")
	fmt.Println(hex.EncodeToString(Mtree.root.right.nodeHash))

	// calculate hash of left and right nodes after appending
	temp := []byte{}
	temp = append(temp, Mtree.root.left.nodeHash...)
	temp = append(temp, Mtree.root.right.nodeHash...)
	// fmt.Println(hex.EncodeToString(temp))
	temp2 := sha256.Sum256(temp)
	// final (supposed to be) merkle hash
	fmt.Println(hex.EncodeToString(temp2[:]))
	fmt.Println(hex.EncodeToString(Mtree.rootHash))
}
