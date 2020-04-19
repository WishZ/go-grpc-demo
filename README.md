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
> 启动服务端
```shell script
go run main.go -grpc-port=9090 -db-host=<HOST>:3306 -db-user=<USER> -db-password=<PASSWORD> -db-schema=<SCHEMA>
```
> 启动客户端
```shell script
go run main.go -server=localhost:9090
```