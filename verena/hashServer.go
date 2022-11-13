package main

// usage:
// To run the server (from command line): go run hashServer.go
// To interact with the server, open a separate terminal:
// To get a value from hashServer: curl -X POST localhost:8090/get -H 'Content-Type: application/json' -d '{"uuid":"VALUE"}'
// To Put a value to hashServer: curl -X POST localhost:8090/put -H 'Content-Type: application/json' -d '{"uuid":"VALUE", "entry":{"hash":"HASH_VALUE", "version":VERSION, "publicKey":"PUBLIC_KEY","hash_old":"HASH_VALUE", "version_old":VERSION, "publicKey_old":"PUBLIC_KEY"}}'

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Entry struct {
	Hash          string `json:"hash"`
	Version       int    `json:"version"`
	PublicKey     string `json:"publicKey"`
	Hash_old      string `json:"hash_old"`
	Version_old   int    `json:"version_old"`
	PublicKey_old string `json:"publicKey_old"`
}

type Value struct {
	Hash      string `json:"hash"`
	Version   int    `json:"version"`
	PublicKey string `json:"publicKey"`
}

func (e Entry) String() string {
	return fmt.Sprintf("[%s, %d, %s, %s, %d, %s]", e.Hash, e.Version, e.PublicKey, e.Hash_old, e.Version_old, e.PublicKey_old)
}

func (e Value) String() string {
	return fmt.Sprintf("[%s, %d, %s]", e.Hash, e.Version, e.PublicKey)
}

type GetRequest struct {
	Uuid string `json:"uuid"`
}

type PutRequest struct {
	Uuid  string `json:"uuid"`
	Entry Entry  `json:"entry"`
}

var kv_store = make(map[string]Value)

func add_value_into_map(e Entry, uuid string) {
	var v Value
	v.Hash = e.Hash
	v.PublicKey = e.PublicKey
	v.Version = e.Version
	kv_store[uuid] = v
	fmt.Println("Successful put for uuid: " + uuid + " (entry: " + v.String() + ")")
}

func get(w http.ResponseWriter, req *http.Request) {
	var jsonData GetRequest
	err := json.NewDecoder(req.Body).Decode(&jsonData)
	if err != nil {
		panic(err)
	}

	uuid := jsonData.Uuid
	kv_value, temp := kv_store[uuid]

	if temp {
		fmt.Fprintf(w, "get succeeded for uuid:"+uuid+"\n"+uuid+" "+kv_value.String())
		fmt.Println("Successful get for uuid: " + uuid)
	} else {
		fmt.Fprintf(w, "get failed:uuid "+uuid+" does not exist")
		fmt.Println("failed get for uuid: " + uuid)
	}
}

func put(w http.ResponseWriter, req *http.Request) {
	var jsonData PutRequest
	err := json.NewDecoder(req.Body).Decode(&jsonData)
	if err != nil {
		panic(err)
	}

	uuid := jsonData.Uuid
	kv_value, temp := kv_store[uuid]

	/* The below code checks
	 * Hash value corresponding to a uuid in the hash_server matches the old_hash value sent by client.
	 * Version corresponding to a uuid in the hash_server matches the old_version value sent by client.
	 * Public key corresponding to a uuid in the hash_server matches the old_publicKey value sent by client.
	 * Client sends a new version (usually an incrementing counter) for every put request .
	 * If the document is saved for the first time, it's version should be set to 1.
	 * Only when the above checks passes the put request is accepted.
	 * Else the put request is rejected inlieu of replay attack.
	 */

	if temp {
		if kv_value.Hash == jsonData.Entry.Hash_old &&
			kv_value.Version == jsonData.Entry.Version_old &&
			kv_value.PublicKey == jsonData.Entry.PublicKey_old {
			if jsonData.Entry.Version_old > kv_value.Version {
				add_value_into_map(jsonData.Entry, uuid)
				fmt.Fprintf(w, uuid+" "+kv_store[uuid].String()+" put succeeded\n")
			} else {
				fmt.Fprintf(w, uuid+" "+jsonData.Entry.String()+" put failed: bad version\n")
				fmt.Println("Failed put for uuid: " + uuid + " (entry: " + jsonData.Entry.String() + ")" + " Bad version\n")
			}
		} else {
			fmt.Fprintf(w, uuid+" "+jsonData.Entry.String()+" put failed: Bad Entry Value \n")
			fmt.Println("Failed put for uuid: " + uuid + " (entry: " + jsonData.Entry.String() + ")" + " Bad Entry Value\n")

		}
	} else {
		if jsonData.Entry.Version == 1 {
			add_value_into_map(jsonData.Entry, uuid)
			fmt.Fprintf(w, uuid+" "+kv_store[uuid].String()+" put succeeded\n")
		} else {
			fmt.Fprintf(w, uuid+" "+jsonData.Entry.String()+" put failed: uuid not present\n")
			fmt.Println("Failed put for uuid: " + uuid + " (entry: " + jsonData.Entry.String() + ")" + " uuid not present\n")
		}
	}
}

func main() {
	// test case
	kv_store["1"] = Value{
		Hash:      "b",
		Version:   1,
		PublicKey: "Alice",
	}

	http.HandleFunc("/get", get)
	http.HandleFunc("/put", put)

	http.ListenAndServe(":8090", nil)
}
