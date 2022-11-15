import requests
import uuid

USERNAME = "Alice"
PASSWORD = "12345"

def create_command(command_name):
    return {
        "username": USERNAME,
        "password": PASSWORD,
        "command": command_name
    }

def load_file(filename):
    r = requests.get("http://localhost:8091/loadFile", json={"username":USERNAME, "password":PASSWORD, "filename":filename})
    print(r.text)

def store_file(filename, content):
    r = requests.put("http://localhost:8091/storeFile", json={"username":USERNAME, "password":PASSWORD, "filename":filename, "content":content})
    print(r.text)

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
