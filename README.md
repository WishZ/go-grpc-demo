# go-grpc-demo
grpc demo

> 生成客户端和服务端代码
```
./third_party/protoc-gen.sh
```

> 文件结构模板
```
https://github.com/golang-standards/project-layout
```
> 启动服务端
```shell script
go run main.go -grpc-port=9090 -http-port=8080 -db-host=<HOST>:3306 -db-user=<USER> -db-password=<PASSWORD> -db-schema=<SCHEMA> -log-level=-1 -log-time-format=2020-04-19T15:04:05.999999999Z07:00
```
> 启动客户端
```shell script
go run main.go -server=localhost:9090
```