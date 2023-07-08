package app

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/itmisx/timer-schedule/config"
	"github.com/itmisx/timer-schedule/config/redis"
	"github.com/itmisx/timer-schedule/config/timer"
	"github.com/itmisx/timer-schedule/internal/define"

	_redis "github.com/go-redis/redis/v8"
	"github.com/itmisx/go-helper"
	"github.com/itmisx/logger"
	"github.com/robfig/cron/v3"
)

// InitSchedule 初始化定时调度
func InitSchedule() {
	cr := cron.New(cron.WithSeconds())
	for _, tc := range config.Config.Timers {
		func(tc timer.TimerConfig) {
			// 任务调度丢失补偿检查
			timerScheduleMissCheck(tc)
			// 创建定时任务
			cr.AddFunc(tc.Spec, func() {
				// 分发定时任务
				dispatchTimerTask(tc)
				// 保存下次调度时间
				saveNextScheduleTime(tc)
			})
		}(tc)
	}
	// 启动定时任务
	cr.Start()
}

// TaskResultCheck 定时任务执行结果检查
func StartTaskResultCheck() {
	tk := time.NewTicker(time.Second)
	defer tk.Stop()
	for {
		<-tk.C
		minScore := "0"
		maxScore := strconv.FormatInt(time.Now().Unix(), 10)
		// 任务执行结果待确认列表
		taskResultConfirmSet, _ := redis.NewDB().ZRevRangeByScore(
			context.Background(),
			define.RedisUnconfirmedTask,
			&_redis.ZRangeBy{
				Min: minScore,
				Max: maxScore,
			},
		).Result()
		// 确认任务执行结果
		if len(taskResultConfirmSet) > 0 {
			var tc timer.TimerConfig
			for _, tcJson := range taskResultConfirmSet {
				err := json.Unmarshal([]byte(tcJson), &tc)
				if err != nil {
					continue
				}
				exists, _ := redis.NewDB().Exists(context.Background(), define.RedisTaskResultPrefix+tc.ID).Result()
				if exists > 0 {
					redis.NewDB().Del(context.Background(), define.RedisTaskResultPrefix+tc.ID)
				} else {
					if tc.RetryCount < tc.MaxRetryTimes {
						logger.Info(
							context.Background(),
							"RetryTimerTask",
							logger.String("taskID", tc.ID),
							logger.Int64("retryCount", tc.RetryCount),
						)
						tc.RetryCount++
						dispatchTimerTask(tc)
					}
				}
			}
		}
		// 移除过期的
		redis.NewDB().ZRemRangeByScore(
			context.Background(),
			define.RedisUnconfirmedTask,
			minScore,
			maxScore,
		)
	}
}

// timerScheduleMissCheck 定时器调度丢失检查
func timerScheduleMissCheck(tc timer.TimerConfig) {
	// 获取下次调度的时间
	next, err := getNextScheduleTime(tc)
	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	// 获取上次定时器计划调度时间
	nextTimeBeforeString, _ := redis.NewDB().HGet(ctx, define.RedisNextScheduleTime, tc.Group+":"+tc.Name).Result()
	if nextTimeBeforeString == "" {
		return
	}

	// 两次计划调度时间不同，则进行补偿
	if nextTimeBefore, err := strconv.ParseInt(nextTimeBeforeString, 10, 64); err != nil {
		if nextTimeBefore != next {
			dispatchTimerTask(tc)
		}
	}
}

// saveNextScheduleTime 保存下次调度的时间
func saveNextScheduleTime(tc timer.TimerConfig) error {
	nextTime, err := getNextScheduleTime(tc)
	if err == nil {
		rdb := redis.NewDB()
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()
		_, err := rdb.HSet(ctx, define.RedisNextScheduleTime, tc.Group+":"+tc.Name, nextTime).Result()
		return err
	}
	return err
}

// GetNextScheduleTime 获取定时任务的下次执行时间
func getNextScheduleTime(tc timer.TimerConfig) (int64, error) {
	// 计算下次的执行
	parser := cron.NewParser(cron.Second |
		cron.SecondOptional |
		cron.Minute |
		cron.Hour |
		cron.Dom |
		cron.Month |
		cron.Dow |
		cron.Descriptor,
	)
	// 计算下次的执行
	schedule, err := parser.Parse(tc.Spec)
	if err != nil {
		// Print error
		logger.Error(context.Background(),
			"schedule parse error",
			logger.String("timerName", tc.Name),
			logger.String("timerGroup", tc.Group),
		)
		return 0, err
	} else {
		return schedule.Next(time.Now()).Unix(), nil
	}
}

// dispatchTimerTask 分发定时器任务
func dispatchTimerTask(tc timer.TimerConfig) {
	rdb := redis.NewDB()
	// 队列的key
	timerTaskList := config.Config.RedisPrefix + tc.Group
	// 移除超出最大队列长度的元素
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	listLength, _ := rdb.LLen(ctx, timerTaskList).Result()
	if listLength > config.Config.MaxTimerTask {
		rdb.RPopCount(ctx, timerTaskList, int(listLength-config.Config.MaxTimerTask))
		logger.Warn(context.Background(),
			"timerTaskList too long",
			logger.String("timerTaskList", timerTaskList),
		)
	}
	// 将定时任务加入到队列中
	if tc.ID == "" {
		tc.ID = helper.RandString(16)
	}
	tcJson, _ := json.Marshal(tc)
	rdb.LPush(ctx, timerTaskList, string(tcJson))

	// 加入待确认任务集合
	{
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()
		if tc.EnableRetry {
			rdb.ZAdd(ctx, define.RedisUnconfirmedTask, &_redis.Z{
				Score:  float64(time.Now().Unix() + tc.RetryDelay),
				Member: string(tcJson),
			})
		}
	}
}
