from hashlib import sha256
import json
import requests
import uuid

USERNAME = "Alice"
PASSWORD = "12345"

class StoreFileResponse:
    def __init__(self, response_text):
        response_json = json.loads(response_text)
        self.hash_root = bytearray(response_json["hashRoot"], "utf-8")
        self.merkle_path = [bytearray(node, "utf-8") for node in response_json["merklePath"].strip(" ").split(" ")]
        self.uuid = response_json["uuid"]
        self.old_entry = Entry(response_json["oldEntry"])

class LoadFileResponse:
    def __init__(self, response_text):
        response_json = json.loads(response_text)
        self.hash_root = bytearray(response_json["hashRoot"], "utf-8")
        self.content = response_json["content"]
        self.entry = Entry(response_json["entry"])

class Entry:
    def __init__(self, response_json):
        self.hash = bytearray(response_json["hash"], "utf-8")
        self.version = response_json["version"]
        self.public_key = response_json["publicKey"]

    def is_valid(self):
        return self.hash and self.version and self.public_key

    def to_json(self):
        return {
            "hash": self.hash.decode("utf-8"),
            "version": self.version,
            "publicKey": self.public_key
        }

    def __str__(self):
        return str(self.to_json())

def load_file(filename):
    r = requests.get("http://localhost:8091/loadFile", json={"username":USERNAME, "password":PASSWORD, "filename":filename})
    response = LoadFileResponse(r.text)
    print("hash root:", response.hash_root)
    print("content:", response.content)
    print("entry:", response.entry)

    if response.hash_root != response.entry.hash:
        print("WARNING: supplied root hash does not match hash server:")
        print(f"    provided from server:      {response.hash_root}")
        print(f"    provided from hash server: {response.entry.hash})")

def store_file(filename, content):
    r = requests.put("http://localhost:8091/storeFile", json={"username":USERNAME, "password":PASSWORD, "filename":filename, "content":content})
    response = StoreFileResponse(r.text)
    print("hash root:", response.hash_root)
    print("merkle path:", response.merkle_path)
    print("old entry:", response.old_entry)

    # inclusion verification
    file_hash = sha256(bytearray(content, "utf-8"))
    print(file_hash.digest())
    print(file_hash.hexdigest())
    root_hash = sha256(response.hash_root)
    for sibling_hash in response.merkle_path[::-1]:
        file_hash.update(sibling_hash)
    if file_hash.digest() != root_hash.digest():
        # failed to verify inclusion proof; bail
        print("hash verification failed")
        print("computed root:", file_hash.digest())
        print("received root:", root_hash.digest())
        # return

    # ask server to ask hash server to write a new entry
    requests.post("http://localhost:8090/put", json={
        "uuid": str(response.uuid),
        "entry": {
            "hash": response.hash_root.decode("utf-8"),
            "version": response.old_entry.version + 1,
            "publicKey": USERNAME
        },
        "oldEntry": response.old_entry.to_json()
    })

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
    print()
    print("store file:")
    store_file("test", "content")
    print()
    print("load file:")
    load_file("test")

    # print()
    # print("hash server test")
    # test_hash_server()
