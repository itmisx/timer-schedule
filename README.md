# Timer-schedule

> 分布式定时任务调度系统

#### 🚀️ 特性

- 定时任务自定义
- 按微服务分组管理
- 避免微服务重复调用
- 支持任务失败重试
- 支持调度丢失补偿

#### 🚀️ 快速开始

- 准备docker环境
- 拉取代码, `git clone https://github.com/itmisx/timer-schedule.git`
- 进入项目目录，`cd timer-schedule`
- 编译，`go build -o build/timer-schedule main.go`
- 打包镜像，`docker build  -t itmisx:timer-schedule . `
- 运行， `docker run -d itmisx:timer-schedule`

#### 🚀️ 配置说明

- redis，redis数据库配置
- redis_prefix, redis key的前缀
- max_timer_task, 每个定时任务的最大数量，自动移除旧的
- 定时器配置,
  - name, 定时器名称
  - group，定时器系统分组
  - spec，cron配置
  - enable_retry，启用失败重试
  - max_retry_times, 最大重试次数
  - retry_delay，检查任务失败的延迟

#### 🚀️ 客户端

- 安装
  `go get -u -v github.com/itmisx/timer-schedule`
- 使用
  ```go
  // 创建定时任务消费客户端
  // 参数1，redis前缀
  // 参数2，系统分组
  // 参数3，redis实例
  timerConsumer:=pkg.New("timer-schedule:","admin",redis.NewDB())
  // 添加指定任务的处理函数
  timerConsumer.AddFunc("timerName1",fn)
  // 启动
  timerConsumer.Start()
  ```
