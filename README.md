# End-to-End Encrypted File Sharing System

This repository contains the code for end-to-end encrypted file sharing system.

For comprehensive documentation, see the Project 2 Spec (https://cs161.org/proj2/start_coding.html).

All functions and implementations in `client/client.go`, and tests in `client_test/client_test.go`.

To understand usages, run `go test -v` inside of the `client_test` directory. Each test case has corresponding explainations and comments. Feel free to write more test case and contribute to it. 

## Client Feature Summerize
1. Authenticate user (`InitUser`, `GetUser`)
2. Save files to server (`User.StoreFile`)
3. Load saved files from the server (`User.LoadFile`)
4. Overwrite saved files on the server 
5. Append the saved files on the server (`User.AppendFile`)
6. Share saved files with other user (`User.CreateInvitation`, `User.AcceptInvitation`)
7. Revoke access to previously shared files (`User.RevokeAccess`)

## Server APIs
1. Keystore is a trusted server to store public keys (`KeystoreSet`, `KeystoreGet`)
2. Datastore is a untrusted server. We use it to store encrypted files (`DatastoreSet`, `DatastoreGet`, `DatastoreDelete`)