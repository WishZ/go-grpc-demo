package cmd

import (
	// mysql驱动
	"context"
	"database/sql"
	"flag"
	"fmt"

	"github.com/WishZ/go-todo-service/pkg/logger"
	"github.com/WishZ/go-todo-service/pkg/protocol/grpc"
	"github.com/WishZ/go-todo-service/pkg/protocol/rest"
	v1 "github.com/WishZ/go-todo-service/pkg/service/v1"
	_ "github.com/go-sql-driver/mysql"
)

type Config struct {
	//gRPC服务启动参数
	//服务监听的TCP端口
	GRPCPort string
	// HTTP/REST网关启动参数部分
	// HTTPPort是通过HTTP/REST网关监听的TCP端口
	//HTTP网关是gRPC服务的包装器
	HTTPPort string

	DataStoreDBHost     string
	DataStoreDBUser     string
	DataStoreDBPassword string
	//数据库名
	DataStoreDBSchema string
	// 日志参数部分
	// LogLevel 是全局日志级别：Debug(-1)，Info(0)，Warn(1)，Error(2)，DPanic(3)，Panic(4)，Fatal(5)
	LogLevel int
	// LogTimeFormat 是记录器的打印时间格式，例如2006-01-02T15:04:05.999999999Z07:00
	LogTimeFormat string
}

func RunServer() error {
	ctx := context.Background()

	var cfg Config
	flag.StringVar(&cfg.GRPCPort, "grpc-port", "", "gRPC port to bind")
	flag.StringVar(&cfg.HTTPPort, "http-port", "", "http port to bind")
	flag.StringVar(&cfg.DataStoreDBHost, "db-host", "", "Database host")
	flag.StringVar(&cfg.DataStoreDBUser, "db-user", "", "Database user")
	flag.StringVar(&cfg.DataStoreDBPassword, "db-password", "", "Database password")
	flag.StringVar(&cfg.DataStoreDBSchema, "db-schema", "", "Database schema")
	flag.IntVar(&cfg.LogLevel, "log-level", -1, "log level")
	flag.StringVar(&cfg.LogTimeFormat, "log-time-format", "2006-01-02T15:04:05.999999999Z07:00", "log time format")
	flag.Parse()

	if len(cfg.GRPCPort) == 0 {
		return fmt.Errorf("invalid TCP port for gRPC server: '%s'", cfg.GRPCPort)
	}

	if len(cfg.HTTPPort) == 0 {
		return fmt.Errorf("invalid TCP port for HTTP gateway: '%s'", cfg.HTTPPort)
	}

	//初始化全局日志记录器
	logger.Init(cfg.LogLevel, cfg.LogTimeFormat)

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
	//运行HTTP网关
	go func() {
		_ = rest.RunServer(ctx, cfg.GRPCPort, cfg.HTTPPort)
	}()
	return grpc.RunServer(ctx, v1Api, cfg.GRPCPort)
}
