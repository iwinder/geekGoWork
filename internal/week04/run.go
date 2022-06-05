package week04

import (
	"context"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"github.com/iwinder/geekGoWork/internal/week04/options"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run() {
	// 日志打印配置
	flag.Set("v", "5")
	flag.Set("logtostderr", "true")
	flag.Parse()
	defer glog.Flush()
	glog.V(3).Infoln(fmt.Sprintf("http server start..."))
	// 初始化配置文件
	globalOption := options.InitConfig()
	// 创建 errgroup
	g, errCxt := errgroup.WithContext(context.Background())
	// 用于手动模拟单个服务错误退出
	serverOut := make(chan struct{})

	r := gin.Default()
	// 创建 server
	server := InitServer(globalOption, r, serverOut)
	// g1 启动 http server
	g.Go(func() error {
		glog.V(3).Infoln(fmt.Sprintf("http server run..."))
		return server.ListenAndServe()
	})

	// GRPC
	userServer := InitUserServer(globalOption.MysqlOption)
	grpcService := grpc.NewServer()

	g.Go(func() error {
		glog.V(3).Infoln(fmt.Sprintf("GRPC server run..."))
		gRpcUrl := InitGRpcServer(grpcService, globalOption, userServer)
		lis, err := net.Listen("tcp", gRpcUrl)
		if err != nil {
			glog.Warning(fmt.Sprintf("\n GRPC Listen error... \n"))
			return err
		}
		return grpcService.Serve(lis)
	})

	//接收信号量，方便结束服务
	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// 用于延迟关闭服务，errgroup.WithContext 使用 context.WithTimeout 还是会被转成 context.WithCancel 无法做到延迟关闭
	graceTimeOut := 30
	timeOut := time.Duration(graceTimeOut) * time.Second
	ctx, cancle := context.WithTimeout(context.Background(), timeOut)
	defer cancle()

	// g2 捕获停止条件，关闭 http 服务
	g.Go(func() error {
		select {
		case <-errCxt.Done():
			glog.Warning(fmt.Sprintf("\n errgroup exit... \n"))
		case osCall := <-stopChan: // 捕获
			glog.Warning(fmt.Sprintf("\n system call: %+v , exit... \n", osCall))
		case <-serverOut: //
			glog.Warning("server will out...")
		}
		grpcService.GracefulStop()
		return server.Shutdown(ctx) // 关闭 http server
	})

	if err := g.Wait(); err != nil {
		glog.Fatalln("\n errgroup error: ", err)
	}
	// 由于 glog.Fatalln 这一步在当前环境实际不会运行
	glog.V(3).Infoln(fmt.Sprintf("all group done!"))

}
