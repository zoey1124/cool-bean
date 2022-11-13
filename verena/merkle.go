/*
Merkle Tree related library.
Util functions influde:
    - GetUUID(username, filename) => UUID
	- VerifyFresh
	-
*/

package main

import (
	"crypto/sha256"
	"reflect"

	mt "github.com/cbergoon/merkletree"
	userlib "github.com/cs161-staff/project2-userlib"
)

/*================================= Util Functions ==================================*/
func ByteLengthNormalize(byteArr []byte, k int) []byte {
	/*
			Return a []byte with length. If input []byte len > k, trim the byte array
		    If input []byte length < k, padding with 0
	*/
	if len(byteArr) >= k {
		return byteArr[:k]
	}
	// Padding array with zero to length of k
	n := len(byteArr)
	for i := 0; i < (k - n); i++ {
		byteArr = append(byteArr, 0)
	}
	return byteArr
}

func GetUUID(username string, filename string) userlib.UUID {
	/*
		Return UUID(H(username||filename))
	*/
	username_byte := ByteLengthNormalize([]byte(username), 16)
	filename_byte := ByteLengthNormalize([]byte(filename), 16)
	UUID, _ := userlib.UUIDFromBytes(userlib.Hash(append(username_byte, filename_byte...)))
	return UUID
}

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

func VerifyFresh(roothash []byte, content Content, sibling_hashes [][]byte) (bool, error) {
	/*
		Recalculate roothash using given merkle path, return true if match with given roothash
	*/
	curr_hash, _ := content.CalculateHash()
	for _, sibling_hash := range sibling_hashes {
		h := sha256.New()
		h.Write(append(curr_hash, sibling_hash...))
		curr_hash = h.Sum(nil)
	}
	return reflect.DeepEqual(curr_hash, roothash), nil
}

func GetMerkleRoot(UUID userlib.UUID) []byte {
	tree_byte, _ := userlib.DatastoreGet(UUID)
	var tree mt.MerkleTree
	userlib.Unmarshal(tree_byte, &tree)
	hashroot := tree.MerkleRoot()
	return hashroot
}

func GetMerklePath() ([][]byte, error) {
	return nil, nil
}
