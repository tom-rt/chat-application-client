# Chat application client

Command line client of a chat application.

## Description

The client is meant to run alongside a [server](https://github.com/tom-rt/chat-application-server).
When launched, the client will try to connect to the server and start a session (see host, port and nickname args below).

## Getting Started

### Dependencies

* Go version 1.17
* [Gorilla websocket](github.com/gorilla/websocket) as the only external dependency.

### Installing

* To install the program:
```
git clone git@github.com:tom-rt/chat-application-client.git
cd chat-application-client
go get .
```

### Executing program

* There is one mandatory argument: nickname, which represents your nickname in the chat room.
* There are two non mandatory arguments: host and port, which represent the server's hostname and port. Their default values are respectively localhost and 8080.
* The client can be started by running main.go file, example:
```
go run main.go -nickname=john
```

* Or by building and executing a binary, example:
```
go build
./client -nickname=john
```

## Author

https://github.com/tom-rt
