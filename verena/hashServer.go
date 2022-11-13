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
)

type Entry struct {
    Hash          string `json:"hash"`
    Version       int    `json:"version"`
    PublicKey     string `json:"publicKey"`
}

func (e Entry) String() string {
    return fmt.Sprintf("[%s, %d, %s]", e.Hash, e.Version, e.PublicKey)
}

type GetRequest struct {
    Uuid string `json:"uuid"`
}

type PutRequest struct {
    Uuid  string `json:"uuid"`
    Entry Entry  `json:"entry"`
    OldEntry Entry `json:"oldEntry"`
}

var kv_store = make(map[string]Entry)

func get(w http.ResponseWriter, req *http.Request) {
    var jsonData GetRequest
    err := json.NewDecoder(req.Body).Decode(&jsonData)
    if err != nil {
        panic(err)
    }

    uuid := jsonData.Uuid
    value, exists := kv_store[uuid]

    if ! exists {
        fmt.Println("no key found: " + uuid)
        fmt.Fprintf(w, "null\n")
        return
    }
    fmt.Fprintf(w, value.String() + "\n")
    fmt.Println("Successful get for uuid: " + uuid)
}

func put(w http.ResponseWriter, req *http.Request) {
    var jsonData PutRequest
    err := json.NewDecoder(req.Body).Decode(&jsonData)
    if err != nil {
        panic(err)
    }

    uuid := jsonData.Uuid
    entry := jsonData.Entry
    oldEntry := jsonData.OldEntry
    value, exists := kv_store[uuid]

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
        fmt.Fprintf(w, "null\n")
        return
    }
    kv_store[uuid] = entry
    fmt.Println("Successful put for uuid: " + uuid + " (entry: " + kv_store[uuid].String() + ")")
    fmt.Fprintf(w, kv_store[uuid].String() + "\n")
}

func main() {
    // test case
    kv_store["1"] = Entry{
        Hash:      "b",
        Version:   1,
        PublicKey: "Alice",
    }

    http.HandleFunc("/get", get)
    http.HandleFunc("/put", put)

    http.ListenAndServe(":8090", nil)
}
