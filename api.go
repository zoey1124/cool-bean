package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "os"
    "github.com/cs161-staff/project2-starter-code/client"
)

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
            client.InitUser(username, password)
        case "GetUser":
            client.GetUser(username, password)
        case "StoreFile":
            user, _ := client.GetUser(username, password)
            filename := input["filename"]
            content := input["content"]
            user.StoreFile(filename, []byte(content))
        case "LoadFile":
            user, _ := client.GetUser(username, password)
            filename := input["filename"]
            user.LoadFile(filename)
        case "AppendFile":
            user, _ := client.GetUser(username, password)
            filename := input["filename"]
            content := input["content"]
            user.AppendToFile(filename, []byte(content))
        case "CreatInvitation":
            user, _ := client.GetUser(username, password)
            filename := input["filename"]
            recipientName := input["recipient_name"]
            user.CreateInvitation(filename, recipientName)
        case "AcceptInvitation":
            user, _ := client.GetUser(username, password)
            senderName := input["sender_name"]
            invitationUUID := input["invitation_uuid"]
            filename := input["filename"]
            user.AcceptInvitation(senderName, invitationUUID, filename)
        case "RevokeAccess":
            user, _ := client.GetUser(username, password)
            filename := input["filename"]
            recipientName := input["recipient_name"]
            user.RevokeAccess(filename, recipientName)
    }
}
