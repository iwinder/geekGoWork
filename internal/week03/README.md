## 第三周作业
> 1. 基于 errgroup 实现一个 http server 的启动和关闭 ，以及 linux signal 信号的注册和处理，要保证能够一个退出，全部注销退出。


## API
| 地址 | 参数       | 说明                  |
|---|----------|---------------------|
| http://127.0.0.1:8088/| user（可选） | 通用路由地址              |
| http://127.0.0.1:8088/healthz| 无        | 健康检查                |
| http://127.0.0.1:8088/shutdown| 无        | 模拟单个 http server 退出 |
 | http://127.0.0.1:8088/timeout| 无        | 模拟延迟20s 展示结果        |

## 目录结构
- [internal/week03/server.go](server.go): 实现具体逻辑
- [cmd/week03/week03.go](../../cmd/week03/week03.go): 实际运行用的 main 方法

## 使用方式

进入[cmd/week03](../../cmd/week03/)，运行里面的 [week03.go](../../cmd/week03/week03.go)即可

## 关键实现

关键实现部分在 [internal/week03/server.go](server.go) 的`RunServer()`

```Golang
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
```
