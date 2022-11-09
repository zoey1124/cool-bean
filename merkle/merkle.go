package main

// playground for merkle for now, change to package later

import (
	"crypto/sha256"
	"log"
	"reflect"

	"github.com/cbergoon/merkletree"
	mt "github.com/cbergoon/merkletree"
)

/*=================== Merkle Tree: Implement the Content Interface ===================*/
type Content struct {
	x []byte // encrypted updated content
}

// CalculateHash hashes the values of a Content
func (t Content) CalculateHash() ([]byte, error) {
	h := sha256.New()
	if _, err := h.Write(t.x); err != nil {
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
	return reflect.DeepEqual(t.x, other.(Content).x), nil
}

/*=====================================================================================*/

/*====================== Merkle Tree Functions for Clients =============================*/
func verifyFresh() {
	// Takes a path and a hashroot as input
	// Return True if node + path => hashroot, False otherwise
}

/*=====================================================================================*/

func main() {
	var list []mt.Content
	list = append(list, Content{x: []byte("A")})
	log.Println(list)
	t, err := merkletree.NewTree(list)
	if err != nil {
		log.Fatal(err)
	}
	// Get the Merkle root of the tree
	mr := t.MerkleRoot()
	log.Println("\n", mr)

	// try how to add one more element and update the merkle tree and hashroot?
	list = append(list, Content{x: []byte("B")})
	err = t.RebuildTreeWith(list)
	if err != nil {
		log.Fatal(err)
	}
	// Get the Merkle root of the tree
	mr = t.MerkleRoot()
	log.Println("\n", mr)
}
