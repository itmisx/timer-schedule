package define

import (
	"timer-schedule/config"
)

var (
	RedisUnconfirmedTask  string
	RedisNextScheduleTime string
	RedisTaskResultPrefix string
)

func InitRedisPrefix() {
	// 定时任务执行结果待确认集合
	RedisUnconfirmedTask = config.Config.RedisPrefix + "unconfirmed-task"
	// 定时器下次调度的时间戳
	RedisNextScheduleTime = config.Config.RedisPrefix + "next-schedule-time"
	// 定时任务执行结果
	RedisTaskResultPrefix = config.Config.RedisPrefix + "task-result:"
}
