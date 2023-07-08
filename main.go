package main

import (
	"github.com/itmisx/timer-schedule/config/boot"
	"github.com/itmisx/timer-schedule/internal/app"
)

func main() {
	// 应用配置初始化
	boot.AppInit()
	// 初始化定时器调度
	app.InitSchedule()
	// 开始任务执行结果检查
	go app.StartTaskResultCheck()
	<-make(chan int)
}
