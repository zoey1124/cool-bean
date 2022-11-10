import os
import json

USERNAME = "Alice"
PASSWORD = "12345"

def create_command(command_name):
    return {
        "username": USERNAME,
        "password": PASSWORD,
        "command": command_name
    }

if __name__ == "__main__":
    client_input = {
        "inputs": [
            create_command("InitUser"),
            create_command("GetUser")
        ]
    }
    with open("input.json", "w") as f:
        f.write(json.dumps(client_input))
    os.system("go run server.go input.json")
