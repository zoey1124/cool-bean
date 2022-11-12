package main

// usage:
// to run the server (from command line): go run hashServer.go
// in a separate terminal:
// to get a value: curl -X POST localhost:8090/get -H 'Content-Type: application/json' -d '{"uuid":"[VALUE]"}'
// curl -X POST localhost:8090/put -H 'Content-Type: application/json' -d '{"uuid":"[VALUE]", "entry":{"hash:":"[HASH_VALUE]", "version":[VERSION], "publicKey:"[PUBLIC_KEY""}}'

import (
    "fmt"
    "encoding/json"
    "net/http"
)

type Entry struct {
    Hash string `json:"hash"`
    Version int `json:"version"`
    PublicKey string `json:"publicKey"`
}

func (e Entry) String() string {
    return fmt.Sprintf("[%s, %d, %s]", e.Hash, e.Version, e.PublicKey)
}

type GetRequest struct {
    Uuid string `json:"uuid"`
}

type PutRequest struct {
    Uuid string `json:"uuid"`
    Entry Entry `json:"entry"`
}

var m = make(map[string]Entry)

func get(w http.ResponseWriter, req *http.Request) {
    var jsonData GetRequest
    err := json.NewDecoder(req.Body).Decode(&jsonData)
    if err != nil {
        panic(err)
    }

    uuid := jsonData.Uuid
    fmt.Fprintf(w, uuid + " " + m[uuid].String() + "\n")

    fmt.Println("Successful get for uuid: " + uuid)
}

func put(w http.ResponseWriter, req *http.Request) {
    var jsonData PutRequest
    err := json.NewDecoder(req.Body).Decode(&jsonData)
    if err != nil {
        panic(err)
    }

    uuid := jsonData.Uuid
    m[uuid] = jsonData.Entry
    fmt.Fprintf(w, uuid + " " + jsonData.Entry.String() + " put\n")

    fmt.Println("Successful put for uuid: " + uuid + " (entry: " + jsonData.Entry.String() + ")")
}

func main() {
    // test case
    m["1"] = Entry{
        Hash: "a",
        Version: 0,
        PublicKey: "Alice",
    }

    http.HandleFunc("/get", get)
    http.HandleFunc("/put", put)

    http.ListenAndServe(":8090", nil)
}
