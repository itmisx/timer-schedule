package timer

// TimerConfig 定时器配置
type TimerConfig struct {
	ID            string `json:"id"`                                             // 定时器id
	Name          string `mapstructure:"name" json:"name"`                       // 定时器名称
	Group         string `mapstructure:"group" json:"group"`                     // 同个系统的定时任务属于一个分组
	Spec          string `mapstructure:"spec" json:"-"`                          // 定时器的规格,精确到秒，参考https://github.com/robfig/cron
	EnableRetry   bool   `mapstructure:"enable_retry" json:"enable_retry"`       // 开启失败重试，默认关闭
	MaxRetryTimes int64  `mapstructure:"max_retry_times" json:"max_retry_times"` // 失败重试的次数
	RetryDelay    int64  `mapstructure:"retry_delay" json:"retry_delay"`         // 失败重试延迟时间，单位秒
	RetryCount    int64  `json:"retry_count"`                                    // 重试次数统计
}
