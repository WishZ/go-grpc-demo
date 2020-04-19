package cmd

import (
	// mysql驱动
	"context"
	"database/sql"
	"flag"
	"fmt"
	"github.com/WishZ/go-grpc-demo/pkg/protocol/grpc"
	v1 "github.com/WishZ/go-grpc-demo/pkg/service/v1"
	_ "github.com/go-sql-driver/mysql"
)

type Config struct {
	//gRPC服务启动参数
	//服务监听的TCP端口
	GRPCPort string

	DataStoreDBHost     string
	DataStoreDBUser     string
	DataStoreDBPassword string
	//数据库名
	DataStoreDBSchema string
}

func RunServer() error {
	ctx := context.Background()

	var cfg Config
	flag.StringVar(&cfg.GRPCPort, "grpc-port", "", "gRPC port to bind")
	flag.StringVar(&cfg.DataStoreDBHost, "db-host", "", "Database host")
	flag.StringVar(&cfg.DataStoreDBUser, "db-user", "", "Database user")
	flag.StringVar(&cfg.DataStoreDBPassword, "db-password", "", "Database password")
	flag.StringVar(&cfg.DataStoreDBSchema, "db-schema", "", "Database schema")
	flag.Parse()

	if len(cfg.GRPCPort) == 0 {
		return fmt.Errorf("invalid TCP port for gRPC server: '%s'", cfg.GRPCPort)
	}
	//添加MySQL驱动程序特定参数来解析 date/time
	param := "parseTime=true"
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?%s", cfg.DataStoreDBUser,
		cfg.DataStoreDBPassword, cfg.DataStoreDBHost, cfg.DataStoreDBSchema, param)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database:%v", err)
	}
	defer db.Close()

	v1Api := v1.NewToDoServiceServer(db)

	return grpc.RunServer(ctx, v1Api, cfg.GRPCPort)
}
