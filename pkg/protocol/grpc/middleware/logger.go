package middleware

import (
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

func codeToLevel(code codes.Code) zapcore.Level {
	if code == codes.OK {
		return zapcore.DebugLevel
	}

	return grpc_zap.DefaultClientCodeToLevel(code)
}

func AddLogging(logger *zap.Logger, opts []grpc.ServerOption) []grpc.ServerOption {
	// 日志记录器的共享选项
	o := []grpc_zap.Option{
		grpc_zap.WithLevels(codeToLevel),
	}
	//确保使用zapLogger记录gPRC库内部的日志语句
	grpc_zap.ReplaceGrpcLoggerV2(logger)
	// 添加一元拦截器
	opts = append(opts, grpc_middleware.WithUnaryServerChain(
		grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
		grpc_zap.UnaryServerInterceptor(logger, o...),
	))
	// 添加流拦截器（此处作为示例添加）
	opts = append(opts, grpc_middleware.WithStreamServerChain(
		grpc_ctxtags.StreamServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
		grpc_zap.StreamServerInterceptor(logger, o...),
	))
	return opts
}
