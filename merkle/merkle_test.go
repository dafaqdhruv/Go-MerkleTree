package merkle

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMerkle(t *testing.T) {

	words := []string{"hello", "this", "is", "a", "merkle", "tree"}
	tree := NewTree(words)

	assert.Equal(t, hex.EncodeToString(tree.RootHash), "c393dac244e2441a310dc3a8ca09b2859e3155e572f5f8074db7f35ec6a5eaaa")

}
