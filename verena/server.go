package main

/*
Usage:
To run server (from Terminal command line): `go run server.go`
To interact with the Server, open a separate terminal:
To load file from Server:
	`curl -X POST localhost:8091/loadFile -H 'Content-Type: application/json' -d '{"username":"<USERNAME>", "password":"<PASSWORD>", "filename":"<FILENAME>"}'`
To store a file to Server:
    `curl -X POST localhost:8091/storeFile -H 'Content-Type: application/json' -d '{"username":"<USERNAME>", "password":"<PASSWORD>", "filename":"<FILENAME>", "content":"<CONTENT>"}'`
*/

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"

	mt "github.com/cbergoon/merkletree"
	"github.com/cs161-staff/project2-starter-code/client"
	userlib "github.com/cs161-staff/project2-userlib"
)

var datastore map[userlib.UUID]FileObject = make(map[userlib.UUID]FileObject)
var DataStore = datastore

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

func GetUUID(username string, filename string) (userlib.UUID, error) {
	/*
		Return UUID(H(username||filename))
	*/
	username_byte := ByteLengthNormalize([]byte(username), 16)
	filename_byte := ByteLengthNormalize([]byte(filename), 16)
	UUID, err := userlib.UUIDFromBytes(userlib.Hash(append(username_byte, filename_byte...)))
	if err != nil {
		return userlib.UUIDNew(), err
	}
	return UUID, nil
}

/*=================== Merkle Tree: Implement the Content Interface ===================*/
type LeafContent struct {
	c []byte // []byte(plaintext_content)
}

// CalculateHash hashes the values of a Content
func (t LeafContent) CalculateHash() ([]byte, error) {
	h := sha256.New()
	if _, err := h.Write(t.c); err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}

// Equals tests for equality of two Contents
func (t LeafContent) Equals(other mt.Content) (bool, error) {
	// DeepEqual returns equal if
	//     1. Both slices are nil or non-nil
	// 	   2. Both slice have the same length
	// 	   3. Corresponding slots have the same value
	return reflect.DeepEqual(t.c, other.(LeafContent).c), nil
}

/*=========================== End of Merkle Tree Implementation ========================*/

type LoadFileRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Filename string `json:"filename"`
}

type LoadFileResponse struct {
	Hashroot string `json:"hashRoot"`
	Content string `json:"content"`
}

type StoreFileRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Filename string `json:"filename"`
	Content  string `json:"content"`
}

type StoreFileResponse struct {
	Hashroot string `json:"hashRoot"`
	MerklePath string `json:"merklePath"`
}

type FileObject struct {
	Plaintext  string         // file content plaintext
	MerkleTree *mt.MerkleTree // merkle tree
	Versions   []mt.Content   // history versions of file content
}

type Entry struct {
	Hash      string `json:"hash"`
	Version   int    `json:"version"`
	PublicKey string `json:"publicKey"`
}

type PutRequest struct {
	UUID     userlib.UUID `json:"uuid"`
	Entry    Entry  `json:"entry"`
	OldEntry Entry  `json:"oldEntry"`
}

/* ============================ API ================================= */

func initUser(username string, password string) {
	client.InitUser(username, password)
}

func getUser(username string, password string) *client.User {
	user, _ := client.GetUser(username, password)
	return user
}

func storeFile(w http.ResponseWriter, req *http.Request) {
	var jsonData StoreFileRequest
	err := json.NewDecoder(req.Body).Decode(&jsonData)
	if err != nil {
		panic(err)
	}
	hashroot, merklePath, err := _storeFile(jsonData.Username, jsonData.Filename, jsonData.Content)
	if err != nil {
		panic(err)
	}
	merklePathString := ""
	for _, h := range merklePath {
		merklePathString += string(h[:]) + " "
	}
	jsonResp, err := json.Marshal(StoreFileResponse{string(hashroot), merklePathString})
    if err != nil {
        panic(err)
    }
    w.Write(jsonResp)
	fmt.Println("Success")
}

func _storeFile(username string, filename string, content string) ([]byte, [][]byte, error) {
	/*
		Return new root hash and sibling node hashes
	*/
	// Get UUID
	UUID, _ := GetUUID(username, filename)
	// Get content and merkle tree
	fileObject, ok := DataStore[UUID]

	if !ok {
		// 1. First time the file has been stored
		leafContent := LeafContent{c: []byte(content)}
		var leaves []mt.Content
		leaves = append(leaves, leafContent)
		merkleTree, err := mt.NewTree(leaves)
		if err != nil {
			return nil, nil, errors.New(strings.ToTitle("Can't build a new merkle tree"))
		}
		fileObject = FileObject{Plaintext: content,
			MerkleTree: merkleTree,
			Versions:   leaves}
		DataStore[UUID] = fileObject
		// Get roothash and merkle path
		roothash := merkleTree.MerkleRoot()
		merklePath, _, err := merkleTree.GetMerklePath(leafContent)
		return roothash, merklePath, nil
	}

	// 2. File existed in DataStore before, update file content
	versions := fileObject.Versions
	new_content := LeafContent{c: []byte(content)}
	versions = append(versions, new_content)
	err := fileObject.MerkleTree.RebuildTreeWith(versions)
	if err != nil {
		return nil, nil, errors.New("Can't rebuild merkle tree with new content")
	}
	fileObject.Versions = versions
	fileObject.Plaintext = content
	DataStore[UUID] = fileObject

	roothash := fileObject.MerkleTree.MerkleRoot()
	merklePath, _, err := fileObject.MerkleTree.GetMerklePath(new_content)
	if err != nil {
		return nil, nil, errors.New("Can't get new merkle path")
	}
	return roothash, merklePath, nil
}

func writeHash(UUID userlib.UUID, entry Entry, oldEntry Entry) {
	/*
		Forward old entry and new entry to Hash Server
	*/
	data, err := json.Marshal(PutRequest{
		UUID:     UUID,
		Entry:    entry,
		OldEntry: oldEntry,
	})
	if err != nil {
		panic(err)
	}
	requestBody := bytes.NewBuffer(data)
	resp, err := http.Post(
		"http://localhost:8090/put",
		"application/json",
		requestBody,
	)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	fmt.Println(string(body[:]))
}

func loadFile(w http.ResponseWriter, req *http.Request) {
	/*
		Print hashroot and file content to client.
	*/
	var jsonData LoadFileRequest
	err := json.NewDecoder(req.Body).Decode(&jsonData)
	if err != nil {
		panic(err)
	}
	hashroot, file_content, err := _loadFile(jsonData.Username, jsonData.Filename)
	if err != nil {
		panic(err)
	}

	jsonResp, err := json.Marshal(LoadFileResponse{string(hashroot), file_content})
    if err != nil {
        panic(err)
    }
    w.Write(jsonResp)
	fmt.Println("Success")
}

func _loadFile(username string, filename string) ([]byte, string, error) {
	// Return: hashroot, file_content
	// 1. get UUID from username and filename
	UUID, err := GetUUID(username, filename)
	if err != nil {
		return nil, "", err
	}
	// 2. get hashroot and content
	var fileObject FileObject
	fileObject = datastore[UUID]
	hashroot := fileObject.MerkleTree.MerkleRoot()
	plaintext := fileObject.Plaintext
	return hashroot, plaintext, nil
}

func appendFile(user *client.User, filename string, content []byte) {
	user.AppendToFile(filename, content)
}

func main() {
	http.HandleFunc("/loadFile", loadFile)
	http.HandleFunc("/storeFile", storeFile)

	http.ListenAndServe(":8091", nil)
	/*
		Example commands in Terminal to run:
		1. store a file
			`curl -X POST localhost:8091/storeFile -H 'Content-Type: application/json' -d '{"username":"Alice", "password":"12345", "filename":"somefile.txt", "content":"This is content"}'`
		2. load a file
			`curl -X POST localhost:8091/loadFile -H 'Content-Type: application/json' -d '{"username":"Alice", "password":"12345", "filename":"somefile.txt"}' --output <FILE>`
	*/
}
