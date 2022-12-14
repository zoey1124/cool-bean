# How to run

##  Start the filesystem server
```
go run server.go
```

## Start the hash server
There is a naive implementation of the hash server in this repository, which you can spin up in a separate terminal:
```
go run hashServer.go
```

Alternatively, you can spin up a [distributed version of the hash server](https://github.com/nivkris/dtrust/blob/dev/hashserver/core-modules/hash_server/README.md) (refer to the README for instructions and an architectural overview).

Either will start an http server (listening on port 8091) that listens for get and put requests to the hash server.

## Start the client:
In a separate terminal:
```
python3 client.py
```

This will create an interactive terminal session waiting for commands. Some example commands:
```
> store <filename> <content>
> load <filename>
> quit
```

# Freshness-guaranteed File Storage Application: Why build it?
In decentralized trust, there is a lot of existing work to guarantee every user gets the correct version of shared data, and how to prevent attackers from compromising shared data. Here, we focus on guaranteeing the freshness of the data. Specifically, our work guarantees that a malicious server cannot serve an outdated version of a document to a client without the client becoming aware of this attack (called a rollback attack).

### Prior Work

In CS161 security class, we've built a file-sharing system that guarantees file confidentiality and integrity when storing in an untrusted storage. Users can store, load, and share files in untrusted storage. We use public-key encryption for confidentiality and digital signiture for integrity. Please refer to [file-storage-app repo](https://github.com/zoey1124/file-storage-app) for source code.

In order to add the freshness guarantee, we need merkle tree data structure. We notice that there are a lot of merkle tree packages existing for most of the languages, so we rely on [cbergoon/merkletree](https://pkg.go.dev/github.com/cbergoon/merkletree) golang package. This package has an A+ go report, and over 400 stars on github.
