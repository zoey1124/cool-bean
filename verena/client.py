import json
import requests
import uuid

USERNAME = "Alice"
PASSWORD = "12345"

class StoreFileResponse:
    def __init__(self, response_text):
        response_json = json.loads(response_text)
        self.hash_root = bytearray(response_json["hashRoot"], "utf-8")
        self.merkle_path = [bytearray(node, "utf-8") for node in response_json["merklePath"].split(" ")]

class LoadFileResponse:
    def __init__(self, response_text):
        response_json = json.loads(response_text)
        self.hash_root = bytearray(response_json["hashRoot"], "utf-8")
        self.content = response_json["content"]

def load_file(filename):
    r = requests.get("http://localhost:8091/loadFile", json={"username":USERNAME, "password":PASSWORD, "filename":filename})
    response = LoadFileResponse(r.text)
    print("hash root:", response.hash_root)
    print("content:", response.content)

def store_file(filename, content):
    r = requests.put("http://localhost:8091/storeFile", json={"username":USERNAME, "password":PASSWORD, "filename":filename, "content":content})
    response = StoreFileResponse(r.text)
    print("hash root:", response.hash_root)
    print("merkle path:", response.merkle_path)

def test_hash_server():
    test_uuid = str(uuid.uuid4())

    print("testing put (no old version)")
    r = requests.post(
        "http://localhost:8090/put",
        json={
            "uuid": test_uuid,
            "entry":{"hash":"a", "version":1, "publicKey":USERNAME},
        })
    print(r.text)

    print("testing get")
    r = requests.get("http://localhost:8090/get", json={"uuid": test_uuid})
    print(r.text)

    print("testing put (updating version)")
    r = requests.post(
        "http://localhost:8090/put",
        json={
            "uuid": test_uuid,
            "entry":{"hash":"c", "version":2, "publicKey":USERNAME},
            "oldEntry":{"hash":"a", "version":1, "publicKey":USERNAME},
        })
    print(r.text)

    print("testing get (after update)")
    r = requests.get("http://localhost:8090/get", json={"uuid": test_uuid})
    print(r.text)

if __name__ == "__main__":
    print("server test")
    store_file("test", "content")
    load_file("test")

    print()
    print("hash server test")
    test_hash_server()
