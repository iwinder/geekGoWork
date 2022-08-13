# 第十三周作业：
## 1. 毕业项目
> 对当下自己项目中的业务，进行一个微服务改造，需要考虑如下技术点：
> 
> 微服务架构（BFF、Service、Admin、Job、Task 分模块）
> API 设计（包括 API 定义、错误码规范、Error 的使用）
> gRPC 的使用
> Go 项目工程化（项目结构、DI、代码分层、ORM 框架）
> 并发的使用（errgroup 的并行链路请求）
> 微服务中间件的使用（ELK、Opentracing、Prometheus、Kafka）
> 缓存的使用优化（一致性处理、Pipeline 优化）

由于时间有限，仅将项目 [qingyucms](https://github.com/iwinder/qingyucms) 完成了部分改造工作,后期会根据实际时间继续完善。

## 毕业总结

本来是在学完云原生课程后想深入了解一些golang的进阶用法以及实际开发中可能涉及到的高阶实践。随着课程的学习，感觉是在上一场golang的微服务架构改造课程，简直是一课两吃。

本门课程对于编程0基础的可能会吃力一些，不过助教大明老师最开始的golang入门以及助教答疑部分还是很好玩很容易理解的，而如果有一定工作经验后，则可以看到以及深入理解各种有趣的架构实践方案。当云原生架构、java微服务、golang微服务三者相互印证，更会发现好多东西是互通的，有的甚至可以相互关联组成知识网络。由于时间问题，目前只能走马观花看一遍，只对事务方面从本地事务->分布式基础原理->分布式事务三个方面整理出三篇笔记，之后会考虑逐步将剩余知识点整理成系列笔记。

课程第一周算是对整个课程方向的总结，之后几周针对知识点开始展开讲解。异常部分讨论了各种语言中对异常的实现方案，之后一步步优化golang的异常实现方案，让我们看到了一个异常部分从诞生到完善的整个思考过程。并发部分主要是goroutine和sync相关锁的讨论，在goroutine中让我们一步步看到调用不规范可能产生的各种内存泄露问题以及优化方案，sync锁部分则说明了并发要考虑的两个问题原子性和可见性。这个是任何语言都无法逃避的，在多语言对比后会发现，他们解决的方式其实也都是差不多的，为了应对重排序等有了 happens-before 规则，在go1.19中内存模型更是和java等对齐。这两个在某种程度上实际也是为了保证最终的数据一致性。如此看来，再扩大到分布式系统中，其各种技术的出现最终也是为了保证一致性。而从这一点上来说，技术又会回到一切的原点计算机系统之上，当前的分布式系统的设计中很多技术以及解决方案都是可以在计算机系统以及网络模型中找到原型的。比如缓存部分，我们常见的有本地缓存、Redis前置缓存，这些可以类比为高速缓存(即三级缓存的总称)与内存。这也是一种一层不够再加一层的实现方案，网络模型的分层也是一种层层叠加，当前的像docker等云原生架构中的overlay方案，也是沿用了这种一层不够再加一层的方案设计，只不过这种反而会增加一定的资源消耗。

像可用性设计中，滑动窗口以及慢启动的设计用来计算吞吐量的方式在 java的 Sentinel 也有涉及到，而这个方案我们则可以追溯到网络中tcp的拥塞控制部分。之后提到的热点数据失效，多个重试节点不断聚拢，最终保证一个执行查询的方案，则可见DNS服务的实现方案。

之后在评论系统的设计架构中则见识到了另一种消息队列的使用方式，对于热点数据的处理则是一个亮点，这里对可能多人触发缓存回填通过结合kafka消息队列回源，产生一个消息到kafka告知某个主题需要更新，job消费后自动慢慢回填，从而实现异步更新缓存的功能。

上面仅是列举的部分，后面还有很多知识点值得参考，课程虽然学完，但还是需要一些时间消化成自己的东西，之后会考虑结合各方面的知识整理成思维导图等自己最终的知识点网络。



