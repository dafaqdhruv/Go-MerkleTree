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

	if n.isLeaf {

		out := sha256.Sum256([]byte(n.val))
		n.nodeHash = out[:]
		fmt.Println(hex.EncodeToString(n.nodeHash) + " for " + n.val)
	} else {

		// if node contains only 1 child
		// transfer the same hash upwards

		if n.right == nil && n.left != nil {
			n.nodeHash = n.left.nodeHash
		} else {

			// else node contains 2 children
			// parent hash =  hash (child1 | child2)
			hashSlice := []byte{}
			hashSlice = append(hashSlice, n.left.nodeHash...)
			hashSlice = append(hashSlice, n.right.nodeHash...)
			fmt.Println(hex.EncodeToString(hashSlice))

			temp := sha256.Sum256(hashSlice)
			n.nodeHash = temp[:]
			fmt.Println(hex.EncodeToString(n.nodeHash))
		}
	}
}

// takes data array input and returns a merkle tree
func buildTree(arr []string) *merkleTree {

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

	for level := math.Ceil(math.Log2(float64(n))) - 1; level >= 0; level-- {

		for nodex := 0; nodex <= int(math.Pow(2, level-1)); nodex += 2 {

			if tempSlice[nodex].nodeHash == nil {
				break
			}
			if !tempSlice[nodex].isLeaf && tempSlice[nodex].right == nil {
				tempParent := &node{
					left:   &tempSlice[nodex],
					right:  nil,
					isLeaf: false,
				}
				tempSlice[nodex].parent = tempParent
				tempParent.getNodeHash()
				tempSlice[nodex/2] = *tempParent

				if level == 1 {
					outTree.root = tempParent
					outTree.rootHash = tempParent.nodeHash
				}
				break
			}

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
}
