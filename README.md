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