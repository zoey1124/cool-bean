package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/cs161-staff/project2-starter-code/client"
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

func loadFile(user *client.User, filename string) {
	user.LoadFile(filename)
}

func appendFile(user *client.User, filename string, content []byte) {
	user.AppendToFile(filename, content)
}

func main() {
	// Read command line argument
	if len(os.Args) < 2 {
		fmt.Println("Please provide input file")
		return
	}
	input_file := os.Args[1]

	// Read input commands from a json file
	jsonFile, err := os.Open(input_file)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Successfully Opened json file")
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var inputs Inputs
	json.Unmarshal([]byte(byteValue), &inputs)

	for i := 0; i < len(inputs.Inputs); i++ {
		input := inputs.Inputs[i]
		command := input.Command
		username := input.Username
		password := input.Password

		// Parse command and call client functions
		fmt.Println(command)
		switch command {
		case "InitUser":
			initUser(username, password)
		case "GetUser":
			getUser(username, password)
		case "StoreFile":
			storeFile(getUser(username, password), input.Filename, []byte(input.Content))
		case "LoadFile":
			loadFile(getUser(username, password), input.Filename)
		case "AppendFile":
			appendFile(getUser(username, password), input.Filename, []byte(input.Content))
		}
	}
}
