package main

// usage:
// to run the server (from command line): go run hashServer.go
// in a separate terminal:
// to get a value: curl -X POST localhost:8090/get -H 'Content-Type: application/json' -d '{"uuid":"[VALUE]"}'
// curl -X POST localhost:8090/put -H 'Content-Type: application/json' -d '{"uuid":"[VALUE]", "hash":"[HASH_VALUE]"}'

import (
    "fmt"
    "encoding/json"
    "net/http"
    // "strconv"
)

var m = make(map[string]string)

func get(w http.ResponseWriter, req *http.Request) {
    decoder := json.NewDecoder(req.Body)
    var jsonData map[string]string
    err := decoder.Decode(&jsonData)
    if err != nil {
        panic(err)
    }
    uuid := jsonData["uuid"]

    fmt.Fprintf(w, uuid + " " + m[uuid] + "\n")
}

func put(w http.ResponseWriter, req *http.Request) {
    decoder := json.NewDecoder(req.Body)
    var jsonData map[string]string
    err := decoder.Decode(&jsonData)
    if err != nil {
        panic(err)
    }
    uuid := jsonData["uuid"]
    value := jsonData["hash"]
    // i, err := strconv.Atoi(value)
    // if err != nil {
    //     panic(err)
    // }
    // fmt.Println(uuid)
    // fmt.Println(i)

    m[uuid] = value
    fmt.Fprintf(w, value + "put\n")
}

func main() {
    m["1"] = "1"

    http.HandleFunc("/get", get)
    http.HandleFunc("/put", put)

    http.ListenAndServe(":8090", nil)
}
