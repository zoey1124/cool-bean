package main

import (
	"github.com/cs161-staff/project2-starter-code/client"
	userlib "github.com/cs161-staff/project2-userlib"
)

type Input struct {
	Command  string `json:"command"`
	Username string `json:"username"`
	Password string `json:"password"`
	Filename string `json:"filename"`
	Content  string `json:"content"`
}

type Inputs struct {
	Inputs []Input `json:"inputs"`
}

func initUser(username string, password string) {
	client.InitUser(username, password)
}

func getUser(username string, password string) *client.User {
	user, _ := client.GetUser(username, password)
	return user
}

func storeFile(user *client.User, filename string, content []byte) {
	user.StoreFile(filename, content)
}

func loadFile(username string, filename string) ([]byte, []byte) {
	// Return: hashroot, file_content
	// a. get UUID from username and filename
	// b. get hashroot
	// c. get file content
	UUID := GetUUID(username, filename)
	hashroot, _ := userlib.DatastoreGet(UUID)
	file_content, _ := userlib.DatastoreGet(UUID)
	return hashroot, file_content
}

func appendFile(user *client.User, filename string, content []byte) {
	user.AppendToFile(filename, content)
}

func main() {
}
