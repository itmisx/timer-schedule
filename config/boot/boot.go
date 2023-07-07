package boot

import (
	"log"

	"timer-schedule/config"
	"timer-schedule/config/redis"
	"timer-schedule/internal/define"

	"github.com/itmisx/logger"
	"github.com/spf13/viper"
)

// AppInit 应用初始化
func AppInit() {
	vp := viper.New()
	vp.AddConfigPath("./etc")
	vp.AddConfigPath("../etc")
	vp.SetConfigName("app")
	vp.SetConfigType("yaml")
	if err := vp.ReadInConfig(); err != nil {
		log.Fatal(err)
	}
	if err := vp.Unmarshal(&config.Config); err != nil {
		log.Fatal(err)
	}
	if config.Config.MaxTimerTask <= 0 {
		config.Config.MaxTimerTask = 100
	}
	// redis初始化
	redis.RedisInit(config.Config.Redis)
	// 初始化前缀定义
	define.InitRedisPrefix()
	// 日志初始化
	logger.Init(config.Config.Log, logger.String("service.name", "timer-schedule"), logger.String("service.version", define.APPVersion))
}
