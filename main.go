package main

import (
	"encoding/hex"
	"fmt"
	"os"

	mt "merkleTree/merkle"
)

func main() {

	Mtree := mt.NewTree(os.Args[1:])

	fmt.Println("Root hash for the input string is ")
	fmt.Println(hex.EncodeToString(Mtree.RootHash))

	Mtree.SaveToJSON()

}
