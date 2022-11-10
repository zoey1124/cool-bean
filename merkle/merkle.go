package main

// playground for merkle for now, change to package later

import (
	"crypto/sha256"
	"fmt"
	"reflect"

	mt "github.com/cbergoon/merkletree"
)

/*=================== Merkle Tree: Implement the Content Interface ===================*/
type Content struct {
	content []byte // content = Encrypted(Compressed(plaintext))
}

// CalculateHash hashes the values of a Content
func (t Content) CalculateHash() ([]byte, error) {
	h := sha256.New()
	if _, err := h.Write(t.content); err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}

// Equals tests for equality of two Contents
func (t Content) Equals(other mt.Content) (bool, error) {
	// DeepEqual returns equal if
	//     1. Both slices are nil or non-nil
	// 	   2. Both slice have the same length
	// 	   3. Corresponding slots have the same value
	return reflect.DeepEqual(t.content, other.(Content).content), nil
}

/*=====================================================================================*/

/*====================== Merkle Tree Functions for Clients =============================*/
// plaintext -> encryption -> Content object
func CovertToContent(plaintext []byte) (Content, error) {
	// enceyption need the secret keys...
	return Content{}, nil
}

// recalculate roothash using given merkle path, return true if match with given roothash
func VerifyFresh(roothash []byte, content Content, sibling_hashes [][]byte) (bool, error) {
	curr_hash, _ := content.CalculateHash()
	for _, sibling_hash := range sibling_hashes {
		h := sha256.New()
		h.Write(append(curr_hash, sibling_hash...))
		curr_hash = h.Sum(nil)
	}
	return reflect.DeepEqual(curr_hash, roothash), nil
}

/*=====================================================================================*/

/*====================== Merkle Tree Functions for Server =============================*/

// Use getMerklePath from library

func main() {
	A := Content{content: []byte("A")}
	B := Content{content: []byte("B")}
	C := Content{content: []byte("C")}
	D := Content{content: []byte("D")}
	content_list := []mt.Content{A, B, C, D}
	tree, _ := mt.NewTree(content_list)
	hashroot := tree.MerkleRoot()
	fmt.Println("Hashroot is", hashroot)
	merkle_path, indexes, _ := tree.GetMerklePath(A)
	fmt.Println("merkle path is ", merkle_path)
	fmt.Println("indexes are", indexes)

	// hash value of content B
	b_hash, _ := B.CalculateHash()
	fmt.Println("B hash is ", b_hash)

	// try to hash(b_hash, b_hash)
	h := sha256.New()
	c_hash, _ := C.CalculateHash()
	d_hash, _ := D.CalculateHash()
	h.Write(append(c_hash, d_hash...))
	fmt.Println(h.Sum(nil))
}
