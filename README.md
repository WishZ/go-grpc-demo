# go-todo-service
grpc demo

> 生成客户端和服务端代码
```
./third_party/protoc-gen.sh
```

> 文件结构模板
```
https://github.com/golang-standards/project-layout
```

## 目录结构
```$xslt
cmd                      主应用目录
internal                 私有应用程序和代码 父级package也只能访问internal package使用大写暴露出的内容，小写的不行
pkg                      共用代码库 业务代码库
vendor                   应用依赖库
api                      OpenAPI / Swagger规范，JSON模式文件，协议定义文件。
web                      静态web资源，服务器端模板和SPA。
configs                  配置文件模板
init                     系统初始化（systemd，upstart，sysv）和进程管理器/主管（runit，supervisor）配置。
scripts                  用于执行各种构建，安装，分析等操作的脚本。
build                    编译后的执行文件
deployments              系统和容器编排部署配置和模板
test                     测试文件
docs                     项目设计文档
tools                    此项目的支持工具。请注意，这些工具可以从/ pkg和/ internal目录导入代码。
examples                 示例
third_party              三方资源
assets                   与资源库一起使用的其他资源(图像，徽标等)
website                  项目站点
githooks                 git钩子
```


> 启动服务端
```shell script
go run cmd/server/main.go -grpc-port=9090 -http-port=8080 -db-host=<HOST>:3306 -db-user=<USER> -db-password=<PASSWORD> -db-schema=<SCHEMA> -log-level=-1 -log-time-format=2006-01-02T15:04:05.999999999Z07:00
```
> 启动客户端
```shell script
#启动 grpc 客户端
go run cmd/client-grpc/main.go -server=localhost:9090
#启动 http 客户端
go run cmd/client-rest/main.go -server=localhost:8080
```