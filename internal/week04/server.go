package week04

import (
	"github.com/gin-gonic/gin"
	api "github.com/iwinder/geekGoWork/api/week04"
	configs "github.com/iwinder/geekGoWork/internal/pkg/options"
	"github.com/iwinder/geekGoWork/internal/week04/service"
	"google.golang.org/grpc"
	"net/http"
)

// InitServer 创建 http server 所需配置
func InitServer(option *configs.Option, g *gin.Engine, serverOut chan struct{}) *http.Server {
	InitRouter(g, serverOut)
	addr := option.ServerOption.HttpServerOption.Host
	if option.ServerOption.HttpServerOption.Port != "" {
		addr = addr + ":" + option.ServerOption.HttpServerOption.Port
	}
	srv := &http.Server{
		Addr:    addr,
		Handler: g,
	}
	return srv
}

// InitGRpcServer gRPC 地址拼接
func InitGRpcServer(grpcService *grpc.Server, option *configs.Option, userServer *service.UserService) string {
	api.RegisterUserServiceServer(grpcService, userServer)
	gRpcUrl := option.ServerOption.GRpcServerOption.Host
	if option.ServerOption.GRpcServerOption.Port != "" {
		gRpcUrl += ":" + option.ServerOption.GRpcServerOption.Port
	}
	return gRpcUrl
}
