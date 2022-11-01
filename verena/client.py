import os
import json

USERNAME = "Alice"
PASSWORD = "12345"

if __name__ == "__main__":
    command = {
        "username": USERNAME,
        "password": PASSWORD,
        "command": "InitUser"
    }
    with open("input.json", "w") as f:
        f.write(json.dumps(command))
    os.system("go run ../api.go")
