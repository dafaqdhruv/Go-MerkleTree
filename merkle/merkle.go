package merkle

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path/filepath"
	"sync"
	"time"

	"oss.terrastruct.com/d2/d2graph"
	elk "oss.terrastruct.com/d2/d2layouts/d2elklayout"
	"oss.terrastruct.com/d2/d2lib"
	"oss.terrastruct.com/d2/d2renderers/d2svg"
	"oss.terrastruct.com/d2/d2themes/d2themescatalog"
	"oss.terrastruct.com/d2/lib/textmeasure"
	"oss.terrastruct.com/util-go/go2"
)

const d2FormatTagToVal = `"%s" -> "%s"
`

const d2FormatTagDeclareHash = `"%s" : {
	"%s"
}
`

type MerkleNode struct {
	parent     *MerkleNode
	LeftChild  *MerkleNode
	RightChild *MerkleNode
	IsLeaf     bool
	Val        string

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

	if n.IsLeaf {
		out := sha256.Sum256([]byte(n.Val))
		n.NodeHash = out[:]
	} else {
		if n.RightChild == nil {
			n.NodeHash = n.LeftChild.NodeHash
		} else {
			// Refer: https://stackoverflow.com/questions/59860517/calculating-sha256-gives-different-results-after-appending-slices-depending-on-i
			var hashSlice []byte
			hashSlice = append(hashSlice, n.LeftChild.NodeHash...)
			hashSlice = append(hashSlice, n.RightChild.NodeHash...)

			// parent hash =  hash (child1 | child2)
			temp := sha256.Sum256(hashSlice)
			n.NodeHash = temp[:]
		}
	}
}

func (n *MerkleNode) BuildTree(arr []string, startIdx int) {

	if len(arr) == 1 {
		n.IsLeaf = true
		n.LeftChild = nil
		n.RightChild = nil
		n.parent = n
		n.Val = arr[0]
		n.Tag = fmt.Sprintf("Hash%d", startIdx)
		n.generateNodeHash()
		return
	}

	pow2 := int(math.Floor(math.Log2(float64(len(arr)))))
	mid := 1 << pow2

	if mid == len(arr) {
		pow2--
		mid = 1 << pow2
	}

	n.Tag = fmt.Sprintf("Hash%d..%d", startIdx, startIdx+len(arr)-1)

	n.LeftChild = NewNode()
	n.LeftChild.parent = n
	n.LeftChild.BuildTree(arr[:mid], startIdx)

	n.RightChild = NewNode()
	n.RightChild.parent = n
	n.RightChild.BuildTree(arr[mid:], startIdx+mid)

	n.generateNodeHash()
}

func (n *MerkleNode) MarshalJSON() ([]byte, error) {
	intermediate := map[string]interface{}{
		"LeftChild":  n.LeftChild,
		"RightChild": n.RightChild,
		"IsLeaf":     n.IsLeaf,
		"Val":        n.Val,
		"NodeHash":   hex.EncodeToString(n.NodeHash),
		"Tag":        n.Tag,
	}
	return json.Marshal(intermediate)
}

func (n *MerkleNode) Walk(fn func(n *MerkleNode, out chan string), out chan string) {
	fn(n, out)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		if n.LeftChild != nil {
			n.LeftChild.Walk(fn, out)
		}
		wg.Done()
	}()

	go func() {
		if n.RightChild != nil {
			n.RightChild.Walk(fn, out)
		}
		wg.Done()
	}()

	wg.Wait()
}

func d2Helper(n *MerkleNode, ch chan string) {
	s := ""

	if n.parent == nil {
		s += fmt.Sprintf(d2FormatTagDeclareHash, n.Tag, hex.EncodeToString(n.NodeHash))
	}

	if n.LeftChild != nil {
		if n.LeftChild.IsLeaf {
			s += fmt.Sprintf(d2FormatTagDeclareHash, n.LeftChild.Tag, n.LeftChild.Val)
		}

		s += fmt.Sprintf(d2FormatTagDeclareHash, n.LeftChild.Tag, hex.EncodeToString(n.LeftChild.NodeHash))
		s += fmt.Sprintf(d2FormatTagToVal, n.Tag, n.LeftChild.Tag)
	}

	if n.RightChild != nil {
		if n.RightChild.IsLeaf {
			s += fmt.Sprintf(d2FormatTagDeclareHash, n.RightChild.Tag, n.RightChild.Val)
		}

		s += fmt.Sprintf(d2FormatTagDeclareHash, n.RightChild.Tag, hex.EncodeToString(n.RightChild.NodeHash))
		s += fmt.Sprintf(d2FormatTagToVal, n.Tag, n.RightChild.Tag)
	}

	ch <- s
}

func (m *MerkleTree) SVGfy() {

	d2Buffer := make(chan string, 2*m.N)
	m.RootNode.Walk(d2Helper, d2Buffer)

	d2tree := func() string {
		out := ""
		for {
			select {
			case f := <-d2Buffer:
				out += f

			case <-time.After(time.Second):
				return out
			}
		}
	}()
	ioutil.WriteFile("out.d2", []byte(d2tree), 0600)

	// d2 SVG render code
	ruler, _ := textmeasure.NewRuler()
	layoutResolver := func(engine string) (d2graph.LayoutGraph, error) {
		return elk.DefaultLayout, nil
	}

	renderOpts := &d2svg.RenderOpts{
		Pad:     go2.Pointer(int64(50)),
		ThemeID: &d2themescatalog.ColorblindClear.ID,
	}

	compileOpts := &d2lib.CompileOptions{
		LayoutResolver: layoutResolver,
		Ruler:          ruler,
	}

	diagram, _, err := d2lib.Compile(context.Background(), d2tree, compileOpts, renderOpts)
	if err != nil {
		log.Panic(err)
	}

	out, err := d2svg.Render(diagram, renderOpts)
	if err != nil {
		log.Panic(err)
	}

	if err := ioutil.WriteFile(filepath.Join("out.svg"), out, 0600); err != nil {
		log.Panic(err)
	}
}

func (m *MerkleTree) SaveToJSON() error {
	bytes, err := json.MarshalIndent(m, "", "\t")
	if err != nil {
		log.Fatal("Error cannot convert tree to JSON", err)
	}

	f, err := os.Create(hex.EncodeToString(m.RootHash))
	if err != nil {
		log.Fatal("Error cannot create output file")
	}
	defer f.Close()

	if _, err := f.Write(bytes); err != nil {
		panic(err)
	}

	return nil
}

func NewNode() *MerkleNode {
	return &MerkleNode{
		IsLeaf:     false,
		LeftChild:  nil,
		RightChild: nil,
		parent:     nil,
		Val:        "",
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

	tree.RootNode.BuildTree(arr, 0)
	tree.RootHash = tree.RootNode.NodeHash

	return tree
}

// To-Do
func (m *MerkleTree) verifyRootHash() {

}
