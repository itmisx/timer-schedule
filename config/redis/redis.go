package redis

import (
	"github.com/itmisx/redisx"
)

var _db redisx.Client

// RedisInit redis初始化
func RedisInit(conf redisx.Config) {
	_db = redisx.New(conf)
}

// NewDB 获取一个新的redis连接
func NewDB() redisx.Client {
	return _db
}
