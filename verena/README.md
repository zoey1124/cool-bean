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

# Freshness-guaranteed File Storage Application: Why build it?

In decentralized trust, there are a lot of existing work with guarantee every user is getting the correct version of a shared data, and how do we prevent from attackers trying to compromise the shared data. This is important for decentralized trust because people need to store file in a untrusted environment, and people need to maintain confedentiality when exchanging information. Here, we are builing on the top of this model: we want to guarantee the freshness of the data. 

### Prior Work 

In CS161 security class, we've built a file-sharing system that guarantees file confidentiality and integrity when storing in an untrusted storage. Users can store, load, and share files in a untrusted storage. We use public-key encryption for confidentiality and digital signiture for integrity. Please refer to [file-storage-app repo](https://github.com/zoey1124/file-storage-app) for source code. 

In order to add the freshness guarantee, we need merkle tree data structure. We notice that there are a lot of merkle traa packages existing for most of the languages, so we rely on [this](https://pkg.go.dev/github.com/cbergoon/merkletree) golang package. 