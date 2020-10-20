## Problem
Write a grpc client talking to the server on grpc stream. 
Client opens the stream with the server and sends an integer counter every 5 secs. 

Server keeps receiving the counter value from the
stream and prints it after reading.

On third time receive from the stream, server closes
the stream with an error.

Client must print the error received from the server
and then reconnect.

After reconnect, client sends the next counter value.

example run:

- client sends '0', server prints 0

- client sends '1', server prints 1

- client sends 2, server prints and closes stream

- client prints error, then reconnects and then send
'3', server prints 3


## Build proto

From root dir run following command

```bash
 protoc -I counter/  -I counter/counter.proto --go_out=plugins=grpc:counter  counter/counter.proto
```


## Server

```bash
go run server/server.go
```

## Client

```bash
go run client/client.go
```