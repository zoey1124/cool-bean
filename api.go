package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "os"
    "github.com/cs161-staff/project2-starter-code/client"
)

func initUser(username string, password string) {
    client.InitUser(username, password)
}

func getUser(username string, password string) (*client.User) {
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
    // Read input commands from a json file
    jsonFile, err := os.Open("input.json")
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println("Successfully Opened json file")
    defer jsonFile.Close()
    byteValue, _ := ioutil.ReadAll(jsonFile)
    var input map[string]string
    json.Unmarshal([]byte(byteValue), &input)
    command := input["command"]
    username := input["username"]
    password := input["password"]

    // Parse command and call client functions
    fmt.Println(command)
    switch command {
        case "InitUser":
            initUser(username, password)
        case "GetUser":
            getUser(username, password)
        case "StoreFile":
            storeFile(getUser(username, password), input["filename"], []byte(input["content"]))
        case "LoadFile":
            loadFile(getUser(username, password), input["filename"])
        case "AppendFile":
            appendFile(getUser(username, password), input["filename"], []byte(input["content"]))
    }
}
