package week03

import (
	"context"
	"encoding/base64"
	"flag"
	"fmt"
	"github.com/golang/glog"
	"golang.org/x/sync/errgroup"
	"io"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func RunServer() {
	// 日志打印配置
	flag.Set("v", "5")
	flag.Set("logtostderr", "true")
	flag.Parse()
	defer glog.Flush()
	glog.V(3).Infoln(fmt.Sprintf("http server start..."))

	// 创建 errgroup
	g, errCxt := errgroup.WithContext(context.Background())
	// 用于手动模拟单个服务错误退出
	serverOut := make(chan struct{})
	// 创建 server
	server := newServer(serverOut)
	// g1 启动 http server
	g.Go(func() error {
		glog.V(3).Infoln(fmt.Sprintf("http server run..."))
		return server.ListenAndServe()
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
		return server.Shutdown(ctx) // 关闭 http server
	})

	if err := g.Wait(); err != nil {
		glog.Warning("\n errgroup error: ", err)
	}
	glog.V(3).Infoln(fmt.Sprintf("all group done!"))
}

// RunServerEror 异常 http server 版本，用作 RunServer 的参照， 这种被 osCall := <-stopChan 阻塞，从而无法执行 g2
// 这种情况将导致调用 /shutdown 将被阻塞
func RunServerEror() {
	flag.Set("v", "5")
	flag.Set("logtostderr", "true")
	flag.Parse()
	defer glog.Flush()

	glog.V(3).Infoln(fmt.Sprintf("http server start..."))
	graceTimeOut := 30
	timeOut := time.Duration(graceTimeOut) * time.Second

	g, errCxt := errgroup.WithContext(context.Background())

	serverOut := make(chan struct{})
	// g1 启动 server
	server := newServer(serverOut)
	g.Go(func() error {
		glog.V(3).Infoln(fmt.Sprintf("http server run..."))
		return server.ListenAndServe()
	})

	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	osCall := <-stopChan // 阻塞后续执行
	glog.V(3).Infoln(fmt.Sprintf("\nsystem call: %+v", osCall))
	ctx, cancle := context.WithTimeout(context.Background(), timeOut)
	defer cancle()

	// g2
	g.Go(func() error {
		select {
		case <-errCxt.Done():
			glog.Warning(fmt.Sprintf("\n errgroup exit... \n"))
		case osCall := <-stopChan:
			glog.Warning(fmt.Sprintf("\n system call: %+v , exit... \n", osCall))
		case <-serverOut:
			glog.Warning("server will out...")
		default:
			glog.Warning(fmt.Sprintf("\n g2 exit... \n"))
		}
		return server.Shutdown(ctx) // 关闭 http server
	})

	if err := g.Wait(); err != nil {
		glog.Warning("\n errgroup error: ", err)
	}
	glog.V(3).Infoln(fmt.Sprintf("all group done!"))
}

// NewServer 创建http server
func newServer(serverOut chan struct{}) *http.Server {
	mux := newServeMux(serverOut)
	return &http.Server{
		Addr:    "127.0.0.1:8088",
		Handler: mux,
	}
}

// NewServeMux 创建路由
func newServeMux(serverOut chan struct{}) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", healthzHandler)
	mux.HandleFunc("/timeout", timeOutHandler)
	mux.HandleFunc("/shutdown", func(w http.ResponseWriter, r *http.Request) {
		glog.V(2).Infoln("Entering shutdown Handler...")
		serverOut <- struct{}{}
	})
	mux.HandleFunc("/", routeHandler)
	return mux
}

// healthzHandler 健康检查路由配置
func healthzHandler(respon http.ResponseWriter, r *http.Request) {
	glog.V(2).Infoln("Entering healthz Handler...")
	io.WriteString(respon, "OK\n")
}

// timeOutHandler 路由配置-延迟返回
func timeOutHandler(reson http.ResponseWriter, r *http.Request) {
	// 业务开始
	glog.V(2).Infoln("Entering timeOut Handler...")
	time.Sleep(20 * time.Second)
	io.WriteString(reson, "Do something 20s OK\n")
}

// routeHandler 通用路由
func routeHandler(reson http.ResponseWriter, r *http.Request) {
	if r.URL.RequestURI() == "/favicon.ico" {
		reson.Header().Set("Content-Type", "image/png")
		base64Str := "data:image/png;base64,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"
		// 用于去掉 “data:image/png;base64,” 部分
		i := strings.Index(base64Str, ",")
		dec := base64.NewDecoder(base64.StdEncoding, strings.NewReader(base64Str[i+1:]))
		io.Copy(reson, dec)
		return
	}
	glog.V(2).Infoln("Entering route Handler...")
	// 获取参数
	user := r.URL.Query().Get("user")
	// 生成随机数
	delay := randInt(0, 2000)
	// 睡眠随机时间，模拟延时
	time.Sleep(time.Millisecond * time.Duration(delay))

	io.WriteString(reson, fmt.Sprintf("hello [%s]\n", user))
}

// randInt 生成min到max的随机数
func randInt(min int, max int) int {
	// rand.Seed通过设置不同值，保证每次执行时生成的随机数不一样
	rand.Seed(time.Now().UTC().UnixNano())
	return min + rand.Intn(max-min)
}
