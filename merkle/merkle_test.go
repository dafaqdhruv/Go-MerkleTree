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

}
