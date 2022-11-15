package main

// usage:
// To run the server (from command line): `go run hashServer.go`
// To interact with the server, open a separate terminal:
// To get a value from hashServer: `curl -X POST localhost:8090/get -H 'Content-Type: application/json' -d '{"uuid":"VALUE"}'`
// To put a value in hashServer:
//     curl -X POST localhost:8090/put -H 'Content-Type: application/json'
//     -d '{"uuid":"VALUE", "entry":{"hash":"HASH_VALUE", "version":VERSION, "publicKey":"PUBLIC_KEY"},
//     "oldEntry":{"hash":"HASH_VALUE", "version":VERSION, "publicKey":"PUBLIC_KEY"}}'

import (
    "encoding/json"
    "fmt"
    "net/http"
    "github.com/google/uuid"
)

type UUID = uuid.UUID

type Entry struct {
    Hash string `json:"hash"`
    Version int `json:"version"`
    PublicKey string `json:"publicKey"`
}

func (e Entry) String() string {
    return fmt.Sprintf("[%s, %d, %s]", e.Hash, e.Version, e.PublicKey)
}

type GetRequest struct {
    UUID UUID `json:"uuid"`
}

type PutRequest struct {
    UUID UUID `json:"uuid"`
    Entry Entry  `json:"entry"`
    OldEntry Entry `json:"oldEntry"`
}

var kv_store = make(map[UUID]Entry)

func writeResponse(w http.ResponseWriter, entry Entry) {
    jsonResp, err := json.Marshal(entry)
    if err != nil {
        panic(err)
    }
    w.Write(jsonResp)
}

func get(w http.ResponseWriter, req *http.Request) {
    var jsonData GetRequest
    err := json.NewDecoder(req.Body).Decode(&jsonData)
    if err != nil {
        panic(err)
    }

    reqUUID := jsonData.UUID
    value, exists := kv_store[reqUUID]

    if ! exists {
        fmt.Println("no key found: " + reqUUID.String())
        writeResponse(w, Entry{})
        return
    }
    writeResponse(w, value)
    fmt.Println("Successful get for uuid: " + reqUUID.String())
}

func put(w http.ResponseWriter, req *http.Request) {
    var jsonData PutRequest
    err := json.NewDecoder(req.Body).Decode(&jsonData)
    if err != nil {
        panic(err)
    }

    reqUUID := jsonData.UUID
    entry := jsonData.Entry
    oldEntry := jsonData.OldEntry
    value, exists := kv_store[reqUUID]

    var errorMessage string
    success := true

    if ! exists {
        if entry.Version != 1 {
            errorMessage = fmt.Sprintf("Invalid version for new entry: %d", entry.Version)
            success = false
        }
    } else {
        if value != oldEntry {
            errorMessage = fmt.Sprintf(
                "Current entry does not match supplied old entry. Current entry: %s, received old entry: %s",
                value.String(),
                oldEntry.String())
            success = false
        } else if entry.Version <= value.Version {
            errorMessage = fmt.Sprintf(
                "Invalid version for update. Current version: %d, received version: %d",
                value.Version,
                entry.Version)
            success = false
        }
    }

    if ! success {
        fmt.Println(errorMessage)
        writeResponse(w, Entry{})
        return
    }
    kv_store[reqUUID] = entry
    newEntry:= kv_store[reqUUID]
    fmt.Println("Successful put for uuid: " + reqUUID.String() + " (entry: " + newEntry.String() + ")")
    writeResponse(w, newEntry)
}

func main() {
    http.HandleFunc("/get", get)
    http.HandleFunc("/put", put)

    http.ListenAndServe(":8090", nil)
}
