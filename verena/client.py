import base64
from hashlib import sha256
import json
import requests
import uuid
import time

USERNAME = "Alice"
PASSWORD = "12345"

class StoreFileResponse:
    def __init__(self, response_text):
        response_json = json.loads(response_text)
        self.root_hash = base64.b64decode(response_json["rootHash"])
        self.merkle_path = [base64.b64decode(node) for node in response_json["merklePath"]]
        self.indexes = response_json["indexes"]
        self.uuid = response_json["uuid"]
        self.old_entry = Entry(response_json["oldEntry"])

class LoadFileResponse:
    def __init__(self, response_text):
        response_json = json.loads(response_text)
        self.root_hash = base64.b64decode(response_json["rootHash"])
        self.content = response_json["content"]
        self.entry = Entry(response_json["entry"])

class Entry:
    def __init__(self, response_json):
        self.hash = base64.b64decode(response_json["hash"]) if response_json["hash"] else None
        self.version = response_json["version"]
        self.public_key = response_json["publicKey"]

    def is_valid(self):
        return self.hash and self.version and self.public_key

    def to_json(self):
        return {
            "hash": base64.b64encode(self.hash).decode() if self.hash else None,
            "version": self.version,
            "publicKey": self.public_key
        }

    def __str__(self):
        return str(self.to_json())

def load_file(filename):
    r = requests.get("http://localhost:8091/loadFile", json={"username":USERNAME, "password":PASSWORD, "filename":filename})
    response = LoadFileResponse(r.text)
    print("content:", response.content)
    print("root hash:", base64.b64encode(response.root_hash).decode())
    # TODO: we need to implement signing to verify this entry is in fact coming from the hash server
    print("entry:", response.entry)

    if response.root_hash != response.entry.hash:
        print("WARNING: supplied root hash does not match hash server:")
        print(f"    provided from server:      {response.root_hash}")
        print(f"    provided from hash server: {response.entry.hash}")
    else:
        print("hash verification passed: file is fresh")

def store_file(filename, content):
    r = requests.put("http://localhost:8091/storeFile", json={"username":USERNAME, "password":PASSWORD, "filename":filename, "content":content})
    response = StoreFileResponse(r.text)
    print("root hash:", list(response.root_hash))
    print("merkle path:", [list(sibling) for sibling in response.merkle_path])
    print("old entry:", response.old_entry)

    # inclusion verification
    print()
    print("verifying inclusion proof")
    file_hash = sha256(bytearray(content, "utf-8"))
    print("file hash:", list(file_hash.digest()))
    root_hash = response.root_hash
    for sibling_hash, index in zip(response.merkle_path, response.indexes):
        if index:
            concat_bytes = file_hash.digest() + sibling_hash
        else:
            concat_bytes = sibling_hash + file_hash.digest()
        file_hash = sha256(concat_bytes)
    if file_hash.digest() != root_hash:
        print("hash verification failed")
        print("computed root:", list(file_hash.digest()))
        print("received root:", list(root_hash))
        return

    # ask server to ask hash server to write a new entry
    print("verification succeeded: asking server to write new entry")
    requests.post("http://localhost:8091/writeHash", json={
        "uuid": str(response.uuid),
        "entry": {
            "hash": base64.b64encode(root_hash).decode(),
            "version": response.old_entry.version + 1,
            "publicKey": USERNAME
        },
        "oldEntry": response.old_entry.to_json()
    })

def test_hash_server():
    test_uuid = str(uuid.uuid4())
    test_hash = base64.b64encode(sha256(bytes("content", "utf-8")).digest()).decode()
    print("randomly generated uuid:", test_uuid)
    print("test hash:", test_hash)
    print()

    print("testing put (no old version)")
    r = requests.post(
        "http://localhost:8090/put",
        json={
            "uuid": test_uuid,
            "entry":{"hash":test_hash, "version":1, "publicKey":USERNAME},
        })
    print("response from hash server:", r.text)
    print()

    print("testing get")
    r = requests.get("http://localhost:8090/get", json={"uuid": test_uuid})
    print("response from hash server:", r.text)
    print()

    print("testing put (updating version)")
    r = requests.post(
        "http://localhost:8090/put",
        json={
            "uuid": test_uuid,
            "entry":{"hash":test_hash, "version":2, "publicKey":USERNAME},
            "oldEntry":{"hash":test_hash, "version":1, "publicKey":USERNAME},
        })
    print("response from hash server:", r.text)
    print()

    print("testing get (after update)")
    r = requests.get("http://localhost:8090/get", json={"uuid": test_uuid})
    print("response from hash server:", r.text)
    print()

def one_store_load(filename, file_content):
    store_file(filename, file_content)
    load_file(filename)
    return

if __name__ == "__main__":
    # do 10 experiments, every experiment 100 store + load
    total_store_time_data = []
    total_load_time_data = []

    # store
    for i in range(1):
        filename = "somefile{}.txt".format(i)
        file_content = "This is a file content"
        start_time = time.time()
        for j in range(100): 
            file_content += str(j)
            store_file(filename, file_content)

        time_used = time.time() - start_time
        total_store_time_data.append(time_used)
    print("time data collected is {}".format(total_store_time_data))
    ave_time = sum(total_store_time_data) / (10)
    print("average time for 100 store is {}".format(ave_time))



    # load
    for i in range(10):
        filename = "somefile{}.txt".format(i)
        start_time = time.time()
        for j in range(100): 
            load_file(filename)

        time_used = time.time() - start_time
        total_load_time_data.append(time_used)
    print("time data collected is {}".format(total_load_time_data))
    ave_time = sum(total_load_time_data) / (10)
    print("average time for 100 load is {}".format(ave_time))
