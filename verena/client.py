import os
import json
import requests

USERNAME = "Alice"
PASSWORD = "12345"

def create_command(command_name):
    return {
        "username": USERNAME,
        "password": PASSWORD,
        "command": command_name
    }

def run_commands(commands):
    with open("input.json", "w") as f:
        f.write(json.dumps(client_input))
    os.system("go run server.go input.json")

def test_hash_server():
    r = requests.get("http://localhost:8090/get", json={"uuid": "1"})
    print(r.text)

    r = requests.post(
        "http://localhost:8090/put",
        json={
            "uuid": "1",
            "entry":{"hash":"a", "version":2, "publicKey":"Alice"},
            "oldEntry":{"hash":"b", "version":1, "publicKey":"Alice"},
        })
    print(r.text)

    r = requests.get("http://localhost:8090/get", json={"uuid": "1"})
    print(r.text)

    r = requests.post(
        "http://localhost:8090/put",
        json={
            "uuid": "2",
            "entry":{"hash":"c", "version":1, "publicKey":"Alice"}
        })
    print(r.text)

    r = requests.get("http://localhost:8090/get", json={"uuid": "2"})
    print(r.text)

if __name__ == "__main__":
    client_input = {
        "inputs": [
            create_command("InitUser"),
            create_command("GetUser")
        ]
    }
    run_commands(client_input)

    test_hash_server()
