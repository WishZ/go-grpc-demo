 protoc -I api/proto/v1 api/proto/v1/todo-service.proto --go_out=plugins=grpc:pkg/api/v1 --proto_path=third_party
 protoc -I api/proto/v1 api/proto/v1/todo-service.proto --proto_path=third_party --grpc-gateway_out=logtostderr=true:pkg/api/v1
 protoc -I api/proto/v1 api/proto/v1/todo-service.proto --proto_path=third_party --swagger_out=logtostderr=true:api/swagger/v1