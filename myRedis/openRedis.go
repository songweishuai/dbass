package myRedis

import (
	"github.com/gomodule/redigo/redis"
	"time"
)

var con redis.Conn


func GetRedisInstance() (redis.Conn, error) {
	if con != nil {
		return con, nil
	}

	con, err := redis.DialTimeout("tcp", "127.0.0.1:6379", time.Second, time.Second, time.Second)
	if err != nil {
		return nil, err
	}

	return con, nil
}

