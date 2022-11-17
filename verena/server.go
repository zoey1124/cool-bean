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
	RootHash []byte `json:"rootHash"`
	Content string `json:"content"`
	Entry Entry `json:"entry"`
}

type StoreFileRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Filename string `json:"filename"`
	Content  string `json:"content"`
}

type StoreFileResponse struct {
	RootHash []byte `json:"rootHash"`
	MerklePath [][]byte `json:"merklePath"`
	Indexes []int64 `json:"indexes"`
	UUID userlib.UUID `json:"uuid"`
	OldEntry Entry `json:"oldEntry"`
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

type GetRequest struct {
	UUID userlib.UUID `json:"uuid"`
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
	uuid, rootHash, merklePath, indexes, err := _storeFile(jsonData.Username, jsonData.Filename, jsonData.Content)
	if err != nil {
		panic(err)
	}
	fmt.Println("root hash:", rootHash)
	fmt.Println("merkle path:", merklePath)
	jsonResp, err := json.Marshal(StoreFileResponse{rootHash, merklePath, indexes, uuid, bytesToEntry(_getHash(uuid))})
    if err != nil {
        panic(err)
    }
    w.Write(jsonResp)
	fmt.Println("Success")
}

func _storeFile(username string, filename string, content string) (userlib.UUID, []byte, [][]byte, []int64, error) {
	/*
		Return new root hash and sibling node hashes
	*/
	// Get UUID
	UUID, _ := GetUUID(username, filename)
	// Get content and merkle tree
	fileObject, ok := DataStore[UUID]

	leafContent := LeafContent{c: []byte(content)}
	leafHash, _ := leafContent.CalculateHash()
	fmt.Println("file hash:", leafHash)
	var leaves []mt.Content
	var merkleTree *mt.MerkleTree

	if ok {
		leaves = fileObject.Versions
	}
	leaves = append(leaves, leafContent)
	merkleTree, err := mt.NewTree(leaves)
	if err != nil {
		return UUID, nil, nil, nil, errors.New(strings.ToTitle("Can't build merkle tree"))
	}
	fileObject = FileObject{
		Plaintext: content,
		MerkleTree: merkleTree,
		Versions: leaves,
	}
	DataStore[UUID] = fileObject

	roothash := fileObject.MerkleTree.MerkleRoot()
	merklePath, indexes, err := fileObject.MerkleTree.GetMerklePath(leafContent)
	if err != nil {
		return UUID, nil, nil, nil, errors.New("Can't get new merkle path")
	}
	return UUID, roothash, merklePath, indexes, nil
}

func loadFile(w http.ResponseWriter, req *http.Request) {
	/*
		Print rootHash and file content to client.
	*/
	var jsonData LoadFileRequest
	err := json.NewDecoder(req.Body).Decode(&jsonData)
	if err != nil {
		panic(err)
	}
	uuid, rootHash, file_content, err := _loadFile(jsonData.Username, jsonData.Filename)
	if err != nil {
		panic(err)
	}

	jsonResp, err := json.Marshal(LoadFileResponse{rootHash, file_content, bytesToEntry(_getHash(uuid))})
    if err != nil {
        panic(err)
    }
    w.Write(jsonResp)
	fmt.Println("Success")
}

func _loadFile(username string, filename string) (userlib.UUID, []byte, string, error) {
	// Return: rootHash, file_content
	// 1. get UUID from username and filename
	UUID, err := GetUUID(username, filename)
	if err != nil {
		return UUID, nil, "", err
	}
	// 2. get rootHash and content
	var fileObject FileObject
	fileObject = datastore[UUID]
	rootHash := fileObject.MerkleTree.MerkleRoot()
	plaintext := fileObject.Plaintext
	return UUID, rootHash, plaintext, nil
}

func appendFile(user *client.User, filename string, content []byte) {
	user.AppendToFile(filename, content)
}

func getHash(w http.ResponseWriter, req *http.Request) {
	var jsonData GetRequest
	err := json.NewDecoder(req.Body).Decode(&jsonData)
	if err != nil {
		panic(err)
	}
	resp := _getHash(jsonData.UUID)
	w.Write(resp)
	fmt.Println(resp)
}

func _getHash(uuid userlib.UUID) []byte {
	data, err := json.Marshal(GetRequest{uuid})
	if err != nil {
		panic(err)
	}
	resp, err := http.Post(
		"http://localhost:8090/get",
		"application/json",
		bytes.NewBuffer(data),
	)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return body
}

func writeHash(w http.ResponseWriter, req *http.Request) {
	var jsonData PutRequest
	err := json.NewDecoder(req.Body).Decode(&jsonData)
	if err != nil {
		panic(err)
	}
	resp := _writeHash(jsonData.UUID, jsonData.Entry, jsonData.OldEntry)
	w.Write(resp)
}

func _writeHash(UUID userlib.UUID, entry Entry, oldEntry Entry) []byte {
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
	if err != nil {
		panic(err)
	}
	return body
}

func bytesToEntry(entryBytes []byte) Entry {
	var entry Entry
	err := json.Unmarshal(entryBytes, &entry)
	if err != nil {
		panic(err)
	}
	return entry
}

func main() {
	http.HandleFunc("/loadFile", loadFile)
	http.HandleFunc("/storeFile", storeFile)
	http.HandleFunc("/getHash", getHash)
	http.HandleFunc("/writeHash", writeHash)

	http.ListenAndServe(":8091", nil)
	/*
		Example commands in Terminal to run:
		1. store a file
			`curl -X POST localhost:8091/storeFile -H 'Content-Type: application/json' -d '{"username":"Alice", "password":"12345", "filename":"somefile.txt", "content":"This is content"}'`
		2. load a file
			`curl -X POST localhost:8091/loadFile -H 'Content-Type: application/json' -d '{"username":"Alice", "password":"12345", "filename":"somefile.txt"}' --output <FILE>`
	*/
}
