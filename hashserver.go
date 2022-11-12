package main
	
import (
    "fmt"
    "net/http"
	"encoding/json"
)

type client_value struct {
	hash_old string
	pubkey_old string
	version_old int
	hash_new string
	pubkey_new string
	version_new int
 }

 type server_value struct {
	hash string
	pubkey string
	version int
 }

var kv_store = make(map[string]server_value)

func parseRequestToJson(req *http.Request) map[string]string {
    decoder := json.NewDecoder(req.Body)
    var jsonData map[string]string
    err := decoder.Decode(&jsonData)
    if err != nil {
        panic(err)
    }
    return jsonData
}

func hs_get(w http.ResponseWriter, req *http.Request){
	jsonData := parseRequestToJson(req)
	key := jsonData["id"]
	var t server_value
	t = kv_store[key]
    fmt.Println(t)
}

func hs_put(w http.ResponseWriter, req *http.Request){
    jsonData1 := parseRequestToJson(req)
	key := jsonData1["id"]
	decoder := json.NewDecoder(req.Body) // adding New: to decode server_value 
    var jsonData map[string]server_value
    err := decoder.Decode(&jsonData)
    if err != nil {
        panic(err)
    }
	var value server_value
	value = jsonData["value"]
	fmt.Println(value)
	kv_store[key] = value
    fmt.Println(kv_store[key])
}

/*
func hs_put(w http.ResponseWriter, req *http.Request) {
    decoder := json.NewDecoder(req.Body)
    var jsonData map[string]client_value
	var jsonData map[string]server_value
    err := decoder.Decode(&jsonData)
    if err != nil {
        panic(err)
    }
	id := jsonData["id"]
	value:= jsonData["client_value"]
	s_value:=json["server_value"]
	if kv_value, exists := kv_store["id"]; exists { //checking for existing ID
	//check new version is greater than old version and if old_hash equals to the hash stored in kv_store
	}	
	if (exists) {
		if value.version_new > kv_value.version && value.hash_old==kv_value.hash{
			s_value.version = value.version_new
			s_value.hash = value.hash_new
			s_value.pubkey = value.pubkey_new
			kv_store["id"] = s_value
			fmt.Fprintf(w, id + " " +kv_store[id] + "Success\n")
		} else {
			fmt.Fprintf(w,"Failed to store, version/hash check failed\n")
		}
	} else { 
		if (value.version_new == 1) { // first time entry for the ID
			s_value.version = value.version_new
			s_value.hash = value.hash_new
			s_value.pubkey = value.pubkey_new
			kv_store["id"] = s_value
			fmt.Fprintf(w, id + " " +kv_store[id] + "Success\n")
		} else {
			fmt.Fprintf(w,"Failed to store, version check failed\n")
		}
	}
}
*/
func main(){
	var s server_value
	s = server_value{"123","456",1}
	kv_store["niv"] = s
	http.HandleFunc("/hs_get", hs_get)
	http.HandleFunc("/hs_put", hs_put)
	http.ListenAndServe(":8090", nil)
}