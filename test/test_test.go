package test

import (
	"fmt"
	"log"
	"testing"
	"time"

	"timer-schedule/config/boot"
	"timer-schedule/config/redis"
	"timer-schedule/internal/app"
	"timer-schedule/pkg"

	"github.com/robfig/cron/v3"
)

func TestSchedule(t *testing.T) {
	// 应用配置初始化
	boot.AppInit()
	// 初始化定时器调度
	app.InitSchedule()
	// 开始任务执行结果检查
	go app.StartTaskResultCheck()
	// group1 消费数据
	group1 := "service1"
	go func() {
		taskConsumer := pkg.New("timer-schedule:", group1, redis.NewDB())
		taskConsumer.AddFunc("test1", func() error {
			log.Println("service1-test1")
			return nil
		})
		taskConsumer.Start()
	}()
	// group2 消费数据
	group2 := "service2"
	go func() {
		taskConsumer := pkg.New("timer-schedule:", group2, redis.NewDB())
		taskConsumer.AddFunc("test2", func() error {
			log.Println("service2-test2")
			return nil
		})
		taskConsumer.Start()
	}()
	<-make(chan int)
}

func TestNext(t *testing.T) {
	parser := cron.NewParser(
		cron.Second |
			cron.SecondOptional |
			cron.Minute |
			cron.Hour |
			cron.Dom |
			cron.Month |
			cron.Dow |
			cron.Descriptor,
	)
	// 计算下次的执行
	schedule, err := parser.Parse("*/2 */2 * * *")
	if err != nil {
		log.Panic(err)
	}
	fmt.Println(schedule.Next(time.Now()))
}
