## 第四周作业
> 1. 按照自己的构想，写一个项目满足基本的目录结构和工程，代码需要包含对数据层、业务层、API 注册，以及 main 函数对于服务的注册和启动，信号处理，使用 Wire 构建依赖。可以使用自己熟悉的框架。

## API
### gRPC
端口 8001 ，服务端相关方法
```Go
// 新增用户 
CreateUser(context.Context, *UserInfo) (*UserInfoReply, error)
// 查询用户
GetUserInfo(context.Context, *UserInfo) (*UserInfoReply, error)
```

### Http
| 地址 | 参数       | 说明                  |
|---|----------|---------------------|
| http://127.0.0.1:8002/| name（可选） | 通用路由地址              |
| http://127.0.0.1:8002/healthz| 无        | 健康检查                |
| http://127.0.0.1:8002/shutdown| 无        | 模拟单个 http server 退出 |
| http://127.0.0.1:8002/timeout| 无        | 模拟延迟20s 展示结果        |
## 关键目录
1. [configs](../../configs) 配置文件
    - 里面包含 启动所需 [week04.yaml](../../configs/week04.yaml) 文件
2. [api/week04](../../api/week04) api 层
   - user.proto 为创建的 proto 文件
   - 在该文件中执行`protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative user.proto` 命令，生成 user.pb.go 和  user_grpc.pb.go 文件
   - 生成文件中的相关用户定义作为 dto 对象
   - 目前只实现了 gRPC 相关接口配置，http 的暂没实现
3. [internal/week04/service](service)  service 层
   - 用于将 dto 与 do 类型互转
   - 最终对 api 层暴露的部分
4. [internal/week04/biz](biz) 业务逻辑层
   - 定义 DO 对象
   - 调用 数据库等实现对数据的处理
   - 通过传递相关数据库操作实例实现对数据层的调用
   - 将 DO 与 PO 类型互转
5. [internal/week04/data](data) 数据层
   - 定义 PO 对象 
   - 完成数据层相关逻辑
6. [internal/week04/write.go](write.go) 依赖注入
   - 通过 wire 工具，手动配置依赖注入所需的实现方法 
   - 在当前目录执行 `wire` 命令，可实现程序调用所需 [internal/week04/wire_gen.go](wire_gen.go) 文件
7. [internal/week04/run.go](run.go) 
   - 主要启动逻辑部分，涉及：
     - 读取配置文件，通过 viper 工具简单读取 yaml文件，暂未涉及通过启动参数指定yaml文件
     - 初始化与启动 基于 gin的 http server
     - 初始化与启动 gRPC 服务
     - 信号处理，根据信号量优雅关闭项目