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

> 文件结构模板
```
https://github.com/golang-standards/project-layout
```