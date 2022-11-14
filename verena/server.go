package main

/*
Usage:
To run server (from Terminal command line): `go run server.go`
To interact with the Server, open a separate terminal:
To load file from Server: `curl -X POST localhost:8090/loadFile -H 'Content-Type: application/json' -d '{"username":"<USERNAME>", "password":"<PASSWORD>", "filename":"<FILENAME>"}'`
To store a file to Server: TODO
*/

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

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
type Content struct {
	content []byte
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

/*=========================== End of Merkle Tree Implementation ========================*/

type LoadFileRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Filename string `json:"filename"`
}

type FileObject struct {
	Content    string
	MerkleTree *mt.MerkleTree
}

func initUser(username string, password string) {
	client.InitUser(username, password)
}

func getUser(username string, password string) *client.User {
	user, _ := client.GetUser(username, password)
	return user
}

func storeFile(username string, filename string, content string) {
	// update merkle tree
	// update file content
}

func loadFile(w http.ResponseWriter, req *http.Request) {
	var jsonData LoadFileRequest
	err := json.NewDecoder(req.Body).Decode(&jsonData)
	if err != nil {
		panic(err)
	}
	hashroot, file_content, err := _loadFile(jsonData.Username, jsonData.Filename)
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(w, "get succedded for hashroot:"+string(hashroot)+"file content:"+file_content)
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
	content := fileObject.Content
	return hashroot, content, nil
}

func appendFile(user *client.User, filename string, content []byte) {
	user.AppendToFile(filename, content)
}

func main() {
	// test case
	UUID, _ := GetUUID("Alice", "somefile")
	someFileContent := []byte("some file content")
	// Generate a Merkle Tree
	var list []mt.Content
	list = append(list, Content{content: someFileContent})
	someMerkleTree, _ := mt.NewTree(list)
	// Generate a FileObject
	fileObject := FileObject{Content: string(someFileContent), MerkleTree: someMerkleTree}
	// Put UUID -> FileObject in DataStore
	DataStore[UUID] = fileObject
	// `curl -X POST localhost:8090/loadFile -H 'Content-Type: application/json' -d '{"username":"Alice", "password":"12345", "filename":"somefile"}'`

	http.HandleFunc("/loadFile", loadFile)

	http.ListenAndServe(":8090", nil)
}
