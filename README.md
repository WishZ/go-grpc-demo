# go-grpc-demo
grpc demo

> 生成客户端和服务端代码
```
protoc -I proc/ proc/search.proto --go_out=plugins=grpc:proc
```

> 服务端
```
protoc --go_out=./ --micro_out=. --proto_path={GOPATH}/pkg/mod --proto_path=./   proto/search.proto
```
