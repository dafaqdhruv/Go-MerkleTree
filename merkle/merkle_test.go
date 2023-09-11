package merkle

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

var sampleInput = []string{"hello", "this", "is", "a", "merkle", "tree"}

func TestMerkleTree(t *testing.T) {

	genericObjs := make([]interface{}, len(sampleInput))
	for i, v := range sampleInput {
		genericObjs[i] = v
	}

	tree := NewTree(genericObjs)

	assert.Equal(t, hex.EncodeToString(tree.RootHash), "c393dac244e2441a310dc3a8ca09b2859e3155e572f5f8074db7f35ec6a5eaaa")
}

func BenchmarkMerkleTreeMult(b *testing.B) {

	x := []int{0}
	genericObjs := make([]interface{}, len(x))
	for i, v := range x {
		genericObjs[i] = v
	}

	b.ResetTimer()
	NewTree(genericObjs)
	for i := 1; i < 1000; i++ {
		genericObjs = append(genericObjs, i)
		NewTree(genericObjs)
	}
}

func TestMerkleProofs(t *testing.T) {
	genericObjs := make([]interface{}, len(sampleInput))
	for i, v := range sampleInput {
		genericObjs[i] = v
	}

	tree := NewTree(genericObjs)

	a, _ := hex.DecodeString("1eb79602411ef02cf6fe117897015fff89f80face4eccd50425c45149b148408")
	b, _ := hex.DecodeString("1bc4a70c8f8296f94c555271a91abe32724d2b9748f9fd8da80337b6cf1270e2")
	c, _ := hex.DecodeString("487e8e3fb58ea5fc6855763fe7a918bda75f564dd0649d8c6b7aefb6f23bd094")
	expectedProof := [][]byte{a, b, c}

	proof, err := tree.GenerateProof("hello")
	assert.Nil(t, err)
	assert.Equal(t, proof.Target, "hello")
	assert.Equal(t, proof.Hashes, expectedProof)
	t.Log("generated proof matches expected proof")
	assert.Equal(t, proof.VerifyProof(), tree.RootHash)

	_, err = tree.GenerateProof("This")
	assert.NotNil(t, err)

	_, err = tree.GenerateProof("does not exists")
	assert.NotNil(t, err)

}
