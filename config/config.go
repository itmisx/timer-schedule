package config

import (
	"timer-schedule/config/timer"

	"github.com/itmisx/logger"
	"github.com/itmisx/redisx"
)

var Config ConfigType

// AppConfig 应用配置
type ConfigType struct {
	Log          logger.Config       `mapstructure:"log"`            // 日志配置
	Redis        redisx.Config       `mapstructure:"redis"`          // redis配置
	Timers       []timer.TimerConfig `mapstructure:"timers"`         // 定时器配置
	RedisPrefix  string              `mapstructure:"redis_prefix"`   // 队列前缀
	MaxTimerTask int64               `mapstructure:"max_timer_task"` // 每隔定时器分组最多的任务限制
}
