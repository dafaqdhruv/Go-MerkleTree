package main

import (
	"encoding/hex"
	"fmt"
	"os"

	mt "merkleTree/merkle"
)

func main() {
	genericObjs := make([]interface{}, len(os.Args[1:]))
	for i, v := range os.Args[1:] {
		genericObjs[i] = v
	}
	Mtree := mt.NewTree(genericObjs)

	fmt.Println("Root hash for the input string is ")
	fmt.Println(hex.EncodeToString(Mtree.RootHash))

	Mtree.SaveToJSON()
	Mtree.SVGfy()
}
