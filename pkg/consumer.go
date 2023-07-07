package pkg

import (
	"context"
	"sync"
	"time"

	"github.com/itmisx/redisx"
	"github.com/tidwall/gjson"
)

type timerConsumer struct {
	redisPrefix string // 定时调度前缀
	groupName   string // 任务分组
	handlerMap  sync.Map
	started     bool
	lock        sync.RWMutex
	rdb         redisx.Client
}

// New 新建消费者
func New(redisPrefix, groupName string, rdb redisx.Client) *timerConsumer {
	return &timerConsumer{
		redisPrefix: redisPrefix,
		groupName:   groupName,
		rdb:         rdb,
	}
}

// AddFunc 增加处理函数
// name 定时器名称
func (c *timerConsumer) AddFunc(name string, handler func() error) {
	c.handlerMap.Store(name, handler)
}

// Start 用协程的方式启动消费，
// 已经启动的则忽略
func (c *timerConsumer) Start() {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.started {
		return
	}
	c.started = true
	go func() {
		for {
			timerItems, _ := c.rdb.BRPop(context.Background(), 0, c.redisPrefix+c.groupName).Result()
			for _, timerItem := range timerItems {
				timerName := gjson.Get(timerItem, "name").String()
				timerEnableRetry := gjson.Get(timerItem, "enable_retry").Bool()
				timerRetryDelay := gjson.Get(timerItem, "retry_delay").Int()
				if timerName == "" {
					continue
				}
				if m, ok := c.handlerMap.Load(timerName); ok {
					if handler, ok := m.(func() error); ok {
						err := handler()
						if err == nil && timerEnableRetry {
							c.rdb.Set(
								context.Background(),
								c.redisPrefix+"task-result:"+gjson.Get(timerItem, "id").String(),
								time.Now().Unix(),
								time.Second*time.Duration(timerRetryDelay+60),
							)
						}
					}
				}
			}
		}
	}()
}
