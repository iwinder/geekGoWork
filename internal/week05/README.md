## 第五周作业
> 1. 参考 Hystrix 实现一个滑动窗口计数器。

## 思路
- 设置 触发熔断的请求总数阈值(不设置时默认值 100，即一个窗口期间的最大请求数) `limitCount` 和 触发熔断的失败率阈值(不设置时默认值 75，即75%)`errorPercentage`
- 设置 滑动窗口的长度(单位ms,不设置时默认1000,即1秒) `timeInMilliseconds` 和  滑动窗口中桶的个数(不设置时默认10，) `limitBucket`。若设置这两个参数需保证能前者整除后者。
- 为实现简单，通过 `timeInMilliseconds/limitBucket` 得到每个桶的窗口时长 `bucketSizeInMs`，之后单独起一个 goroutine 每隔 `bucketSizeInMs` 新增一个桶
  - 若当前时间超过窗口期，重新生成桶列表（通过切片），并新增一个桶并放入桶列表
  - 若当前时间在窗口期，直接新增一个桶并放入桶列表。若通列表长度超过 `limitBucket` 限定长度，抛弃最开始的一个，即取切片`[1:]`
- 进入 api 之前校验当前是否在熔断期间，如果是直接返回异常。反之判等当前总请求数是否大于等于 `limitCount` 或者总请求数/失败总请求数是否大于 `errorPercentage`，满足其一说明需要熔断，此时回到请求返回异常
- 若不再熔断期或者无需熔断，继续执行下面的请求操作，最后根据操作成功与否，增加请求数(也可以考虑先增加总请求数，若失败再增加失败请求数,这里为了简单故一起增加)

## 实现说明
- [bucket.go](bucket.go) 中是关于桶的定义
- [rolling.go](rolling.go) 中是关于滑动窗口的定义以及具体实现
  - `NewRollingWindow` 新建实例
  - `RunWindow` 用于定时创建最新的桶， 需放到 goroutine 中
  - `CheckBroken` 用于检测是否在熔断期间，放到过滤器/拦截器中使用
  - `RecordReqResult`用于新增请求，传递 `false`时表示这是个异常请求
  - `ShowAllBucket` 用于展示当前所有桶中的内容，可以和 `CheckBroken`  放在一起
  - `ShutDone` 用于结束滑动窗口的 goroutine，可放在httpServer结束的部分
- [server.go](server.go) 基于 [week03](../week03/) 改造的服务实现，增加了 `rolling` 的引入，以及`healthzHandler`请求拦截器，用于测试效果
- [cmd/week05/week05.go](../../cmd/week05/week05.go) 实际执行命令，可用于运行测试效果

