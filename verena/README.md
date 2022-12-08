# How to run
In [dtrust](https://github.com/dtrust-project/dtrust), 

Spin up the dtrust nodes: `./start_servers.sh`

Start up the hash server: `go run hashServer.go`

This will start an http server (listening on port 8091) that listens for get and put requests to the hash server.

In verena directory, run `go run server.go`. 
Then in a separate terminal, run `python3 client.py`. Then you will see an interactive terminal session waiting for commands. Some example commands can be 
```
> store <filename> <content>
> load <filename>
> quit
```