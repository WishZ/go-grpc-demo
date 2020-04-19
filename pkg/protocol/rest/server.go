package rest

import (
	"context"
	v1 "github.com/WishZ/go-grpc-demo/pkg/api/v1"
	"github.com/WishZ/go-grpc-demo/pkg/logger"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net/http"
	"os"
	"os/signal"
	"time"
)

//允许HTTP / REST网关
func RunServer(ctx context.Context, grpcPort, httpPort string) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	if err := v1.RegisterToDoServiceHandlerFromEndpoint(ctx, mux, "localhost:"+grpcPort, opts); err != nil {
		logger.Log.Fatal("failed to start HTTP gateway: %v\n", zap.String("reason", err.Error()))
	}

	srv := &http.Server{
		Addr:    ":" + httpPort,
		Handler: mux,
	}

	//关闭
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			logger.Log.Warn("shutting down HTTP/REST gateway server...")
			<-ctx.Done()
		}

		_, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		_ = srv.Shutdown(ctx)
	}()
	logger.Log.Info("starting HTTP/REST gateway...")
	return srv.ListenAndServe()
}
