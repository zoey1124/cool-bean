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
        self.content = response_json["content"]

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


def store_file(filename, content):
    r = requests.put("http://localhost:8091/storeFile", json={"username":USERNAME, "password":PASSWORD, "filename":filename, "content":content})

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
    for i in range(10):
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