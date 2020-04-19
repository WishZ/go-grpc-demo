package grpc

import (
	"context"
	v1 "github.com/WishZ/go-grpc-demo/pkg/api/v1"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
)

func RunServer(ctx context.Context, v1API v1.ToDoServiceServer, port string) error {
	listen, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}
	//服务注册
	server := grpc.NewServer()
	v1.RegisterToDoServiceServer(server, v1API)

	//服务关闭
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			//信号是CTRL+C
			log.Println("shutting down gRPC server...")
			server.GracefulStop()
			<-ctx.Done()
		}
	}()

	//启动gRPC服务

	log.Println("starting gRPC server...")

	if err := server.Serve(listen); err != nil {
		log.Fatal("starting gRPC server failed...")
		return err
	}

	return nil
}
